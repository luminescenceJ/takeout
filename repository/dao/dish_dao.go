package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"takeout/common"
	"takeout/internal/api/admin/request"
	"takeout/internal/api/admin/response"
	"takeout/internal/model"
	"takeout/repository"
)

type DishDao struct {
	db *gorm.DB
}

// Transaction 开始一个事务
func (d DishDao) Transaction(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx).Begin()
}

// Insert 插入菜品
func (d DishDao) Insert(transaction *gorm.DB, dish *model.Dish) error {
	if err := transaction.Create(dish).Error; err != nil {
		return fmt.Errorf("failed to insert dish: %w", err)
	}
	return nil
}

// PageQuery 分页查询菜品
func (d DishDao) PageQuery(ctx context.Context, dto *request.DishPageQueryDTO) (*common.PageResult, error) {
	var (
		pageResult common.PageResult
		dishList   []response.DishPageVo
		err        error
	)

	// 构建查询
	query := d.db.WithContext(ctx).Model(model.Dish{})

	// 根据条件过滤查询
	if dto.Name != "" {
		query = query.Where("dish.name LIKE ?", "%"+dto.Name+"%")
	}
	if dto.Status != "" {
		query = query.Where("dish.status = ?", dto.Status)
	}
	if dto.CategoryId != "" {
		query = query.Where("dish.category_id = ?", dto.CategoryId)
	}

	// 查询总记录数
	if err = query.Count(&pageResult.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count dishes: %w", err)
	}

	// 分页查询菜品
	if err = query.Scopes(pageResult.Paginate(&dto.Page, &dto.PageSize)).
		Select("dish.*, c.name as category_name").
		Joins("LEFT OUTER JOIN category c ON c.id = dish.category_id").
		Order("dish.create_time desc").
		Scan(&dishList).Error; err != nil {
		return nil, fmt.Errorf("failed to query dishes: %w", err)
	}

	// 返回分页结果
	pageResult.Records = dishList
	return &pageResult, nil
}

// GetById 根据ID获取菜品
func (d DishDao) GetById(ctx context.Context, id uint64) (*model.Dish, error) {
	var dish model.Dish
	err := d.db.WithContext(ctx).Preload("Flavors").Find(&dish, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get dish by ID: %w", err)
	}
	return &dish, nil
}

// List 根据分类ID查询菜品列表
func (d DishDao) List(ctx context.Context, categoryId uint64) ([]model.Dish, error) {
	var dishes []model.Dish
	err := d.db.WithContext(ctx).Where("category_id = ?", categoryId).Find(&dishes).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list dishes: %w", err)
	}
	return dishes, nil
}

// OnOrClose 启用或关闭菜品
func (d DishDao) OnOrClose(ctx context.Context, id uint64, status int) error {
	if err := d.db.WithContext(ctx).Model(&model.Dish{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update dish status: %w", err)
	}
	return nil
}

// Update 更新菜品
func (d DishDao) Update(transaction *gorm.DB, dish model.Dish) error {
	if err := transaction.Updates(&dish).Error; err != nil {
		return fmt.Errorf("failed to update dish: %w", err)
	}
	return nil
}

// Delete 删除菜品
func (d DishDao) Delete(transaction *gorm.DB, id uint64) error {
	if err := transaction.Model(&model.Dish{}).Where("id = ?", id).Delete(&model.Dish{}).Error; err != nil {
		return fmt.Errorf("failed to delete dish: %w", err)
	}
	return nil
}

// NewDishRepo 新建DishRepo
func NewDishRepo(db *gorm.DB) repository.DishRepo {
	return &DishDao{db: db}
}

//package dao
//
//import (
//	"context"
//	"gorm.io/gorm"
//	"takeout/common"
//	"takeout/internal/api/admin/request"
//	"takeout/internal/api/admin/response"
//	"takeout/internal/model"
//	"takeout/repository"
//)
//
//type DishDao struct {
//	db *gorm.DB
//}
//
//func (d DishDao) Transaction(ctx context.Context) *gorm.DB {
//	return d.db.WithContext(ctx).Begin()
//}
//
//func (d DishDao) Insert(transaction *gorm.DB, dish *model.Dish) error {
//	return transaction.Create(dish).Error
//}
//
//func (d DishDao) PageQuery(ctx context.Context, dto *request.DishPageQueryDTO) (*common.PageResult, error) {
//	var (
//		pageResult common.PageResult
//		dishList   []response.DishPageVo
//		err        error
//	)
//
//	query := d.db.WithContext(ctx).Model(model.Dish{})
//
//	if dto.Name != "" {
//		query = query.Where("dish.name LIKE ?", "%"+dto.Name+"%")
//	}
//	if dto.Status != "" {
//		query = query.Where("dish.status = ?", dto.Status)
//	}
//	if dto.CategoryId != "" {
//		query = query.Where("dish.category_id = ?", dto.CategoryId)
//	}
//
//	if err = query.Count(&pageResult.Total).Error; err != nil {
//		return nil, err
//	}
//
//	// 3.通用分页查询
//	if err = query.Scopes(pageResult.Paginate(&dto.Page, &dto.PageSize)).
//		Select("dish.*,c.name as category_name").
//		Joins("LEFT OUTER JOIN category c ON c.id = dish.category_id").
//		Order("dish.create_time desc").Scan(&dishList).Error; err != nil {
//		return nil, err
//	}
//
//	pageResult.Records = dishList
//	return &pageResult, err
//}
//
//func (d DishDao) GetById(ctx context.Context, id uint64) (*model.Dish, error) {
//	var (
//		dish model.Dish
//		err  error
//	)
//	dish.Id = id
//	err = d.db.WithContext(ctx).Preload("Flavors").Find(&dish).Error
//	// // 逻辑外键 两次查询
//	//if err = d.db.WithContext(ctx).First(&dish).Error; err != nil {
//	//	return nil, err
//	//}
//	//err = d.db.WithContext(ctx).Where("dish_id = ?", dish.Id).Find(&dish.Flavors).Error
//	return &dish, err
//}
//
//func (d DishDao) List(ctx context.Context, categoryId uint64) ([]model.Dish, error) {
//	res := []model.Dish{}
//	err := d.db.WithContext(ctx).Where("category_id = ?", categoryId).Find(&res).Error
//	return res, err
//}
//
//func (d DishDao) OnOrClose(ctx context.Context, id uint64, status int) error {
//	return d.db.WithContext(ctx).Model(&model.Dish{}).Where("id = ?", id).Update("status", status).Error
//}
//
//func (d DishDao) Update(transaction *gorm.DB, dish model.Dish) error {
//	if err := transaction.Updates(&dish).Error; err != nil {
//		return err
//	}
//	// // 更新菜品的另一种方法 : 逐个更新 DishFlavor
//	//for _, flavor := range dish.Flavors {
//	//	if err := db.Updates(&flavor).Error; err != nil {
//	//		return err
//	//	}
//	//}
//	return nil
//}
//
//func (d DishDao) Delete(transaction *gorm.DB, id uint64) error {
//	return transaction.Model(&model.Dish{}).Where("id = ?", id).Delete(&model.Dish{}).Error
//}
//
//func NewDishRepo(db *gorm.DB) repository.DishRepo {
//	return &DishDao{db: db}
//}
