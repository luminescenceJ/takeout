package dao

import (
	"context"
	"errors"
	"github.com/ulule/deepcopier"
	"gorm.io/gorm"
	"takeout/global"
	"takeout/internal/api/user/request"
	"takeout/internal/model"
	"takeout/repository"
)

type ShoppingCartDao struct {
	db         *gorm.DB
	dishDao    DishDao
	setmealDao SetMealDao
}

// Add 添加购物车
func (s ShoppingCartDao) Add(ctx context.Context, shoppingCart model.ShoppingCart) error {
	// 先判断该商品是否存在于购物车
	var (
		shoppingCartList []model.ShoppingCart
		err              error
	)
	shoppingCartList = s.queryShoppingCart(ctx, shoppingCart)
	if shoppingCartList != nil && len(shoppingCartList) == 1 {
		shoppingCart = shoppingCartList[0]
		shoppingCart.Number++
		if err = s.db.Updates(&shoppingCart).Error; err != nil {
			global.Log.Warn("更新用户购物车数量失败！")
			return err
		}
	} else {
		// 如果不存在，插入数据，数量就是1
		// 判断当前添加到购物车的是菜品还是套餐
		if shoppingCart.DishId != 0 {
			dish, err := s.dishDao.GetById(ctx, uint64(shoppingCart.DishId))
			if err != nil {
				return err
			}
			shoppingCart.Name = dish.Name
			shoppingCart.Image = dish.Image
			shoppingCart.Amount = dish.Price
		} else {
			// 添加到购物车的是套餐
			// 查询添加的套餐信息
			setmeal, err := s.setmealDao.GetByIdWithDish(s.setmealDao.db, uint64(shoppingCart.SetmealId))
			if err != nil {
				return err
			}
			shoppingCart.Name = setmeal.Name
			shoppingCart.Image = setmeal.Image
			shoppingCart.Amount = setmeal.Price
		}
		shoppingCart.Number = 1
		if err = s.db.Create(&shoppingCart).Error; err != nil {
			return err
		}
	}
	return nil
}

// List 查看购物车所有
func (s ShoppingCartDao) List(ctx context.Context, userId int) ([]model.ShoppingCart, error) {
	var shoppingCartList []model.ShoppingCart
	shoppingCartList = s.queryShoppingCart(ctx, model.ShoppingCart{UserId: userId})
	return shoppingCartList, nil
}

// 查询购物车数据
func (s ShoppingCartDao) queryShoppingCart(ctx context.Context, shoppingCart model.ShoppingCart) []model.ShoppingCart {
	var shoppingCartList []model.ShoppingCart
	if err := s.db.WithContext(ctx).
		Where(&shoppingCart).
		Order("create_time desc").
		Find(&shoppingCartList).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return shoppingCartList
}

// Clean 清空购物车
func (s ShoppingCartDao) Clean(ctx context.Context, userId int) error {
	if err := s.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&model.ShoppingCart{}).Error; err != nil {
		return err
	}
	return nil
}

// Subtract 减少购物车
func (s ShoppingCartDao) Subtract(ctx context.Context, ShoppingCartDTO request.ShoppingCartDTO, userId int) error {
	var (
		shoppingCart model.ShoppingCart
		err          error
	)
	if err = deepcopier.Copy(ShoppingCartDTO).To(&shoppingCart); err != nil {
		return err
	}
	shoppingCart.UserId = userId
	shoppingCartList := s.queryShoppingCart(ctx, shoppingCart)
	if shoppingCartList == nil || len(shoppingCartList) == 0 {
		global.Log.Warn("错误的购物车减少")
		return nil
	} else {
		shoppingCart = shoppingCartList[0]
		if shoppingCart.Number == 1 {
			err = s.deleteShoppingCartById(ctx, shoppingCart)
			if err != nil {
				return err
			}
		} else {
			shoppingCart.Number--
			if err = s.db.WithContext(ctx).Updates(&shoppingCart).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// deleteShoppingCartById 根据id删除购物车
func (s ShoppingCartDao) deleteShoppingCartById(ctx context.Context, shoppingCart model.ShoppingCart) error {
	if err := s.db.WithContext(ctx).Delete(&shoppingCart).Error; err != nil {
		return err
	}
	return nil
}

// InsertBatchShoppingCart 批量插入购物车数据
func (s ShoppingCartDao) InsertBatchShoppingCart(ctx context.Context, shoppingCartList []model.ShoppingCart) error {
	if err := s.db.WithContext(ctx).Create(&shoppingCartList).Error; err != nil {
		return err
	}
	return nil
}

func NewShoppingCartDao(db *gorm.DB) repository.ShoppingCartRepo {
	return &ShoppingCartDao{
		db:         db,
		dishDao:    DishDao{db: db},
		setmealDao: SetMealDao{db: db},
	}
}
