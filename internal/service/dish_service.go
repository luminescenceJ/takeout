package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"takeout/common"
	"takeout/common/enum"
	"takeout/common/utils"
	"takeout/global"
	"takeout/internal/api/admin/request"
	"takeout/internal/api/admin/response"
	"takeout/internal/model"
	"takeout/repository"
)

// DishCacheKey redis key 菜品缓存key
const DishCacheKey = "dishCache::"

type IDishService interface {
	AddDishWithFlavors(ctx context.Context, dto request.DishDTO) error
	PageQuery(ctx context.Context, dto *request.DishPageQueryDTO) (*common.PageResult, error)
	GetByIdWithFlavors(ctx context.Context, id uint64) (response.DishVo, error)
	List(ctx context.Context, categoryId uint64) ([]response.DishListVo, error)
	OnOrClose(ctx context.Context, id uint64, status int) error
	Update(ctx context.Context, dto request.DishUpdateDTO) error
	Delete(ctx context.Context, ids string) error
}

type DishServiceImpl struct {
	repo           repository.DishRepo
	dishFlavorRepo repository.DishFlavorRepo
}

func (d DishServiceImpl) AddDishWithFlavors(ctx context.Context, dto request.DishDTO) error {
	var err error
	price, _ := strconv.ParseFloat(dto.Price, 64)

	transaction := d.repo.Transaction(ctx)
	defer func() {
		if r := recover(); r != nil {
			// 发生 panic 时回滚事务
			transaction.Rollback()
		} else if err != nil {
			// 发生错误时回滚事务
			transaction.Rollback()
		}
	}()

	dish := model.Dish{
		Id:          0,
		Name:        dto.Name,
		CategoryId:  dto.CategoryId,
		Price:       price,
		Image:       dto.Image,
		Description: dto.Description,
		Status:      enum.ENABLE,
	}
	if err = d.repo.Insert(transaction, &dish); err != nil {
		return err
	}

	// 外键设置
	for i := range dto.Flavors {
		dto.Flavors[i].DishId = dish.Id
	}
	if err = d.dishFlavorRepo.InsertBatch(transaction, dto.Flavors); err != nil {
		return err
	}

	// 最终返回时，提交事务，若提交失败，返回错误
	if err = transaction.Commit().Error; err != nil {
		return err // 这里会直接返回错误，defer 中的回滚会执行一次
	}
	utils.CleanCache(DishCacheKey + "*")
	return nil
}

func (d DishServiceImpl) PageQuery(ctx context.Context, dto *request.DishPageQueryDTO) (*common.PageResult, error) {

	return d.repo.PageQuery(ctx, dto)
}

func (d DishServiceImpl) GetByIdWithFlavors(ctx context.Context, id uint64) (response.DishVo, error) {
	var (
		dish   *model.Dish
		dishVO response.DishVo
		err    error
	)

	if dish, err = d.repo.GetById(ctx, id); err != nil {
		return response.DishVo{}, err
	}

	dishVO = response.DishVo{
		CategoryId:  dish.CategoryId,
		Id:          dish.Id,
		Name:        dish.Name,
		Description: dish.Description,
		Price:       dish.Price,
		Image:       dish.Image,
		Status:      dish.Status,
		Flavors:     dish.Flavors,
		UpdateTime:  dish.UpdateTime,
	}

	return dishVO, err
}

func (d DishServiceImpl) List(ctx context.Context, categoryId uint64) ([]response.DishListVo, error) {
	var (
		dishes    []model.Dish
		cacheData string
		dishList  []response.DishListVo
		err       error
	)
	// 优先命中Redis缓存
	cacheData, err = global.RedisClient.Get(DishCacheKey + strconv.Itoa(int(categoryId))).Result()
	if err == nil {
		global.Log.Info("查询Redis菜品缓存数据: " + DishCacheKey + strconv.Itoa(int(categoryId)))
		if err = json.Unmarshal([]byte(cacheData), &dishList); err == nil {
			return dishList, nil
		}
	} else {
		global.Log.Info("查询Redis菜品缓存数据失败")
	}

	dishes, err = d.repo.List(ctx, categoryId)
	if err != nil {
		return nil, err
	}
	count := len(dishes)
	dishVo := make([]response.DishListVo, count)
	for i := 0; i < count; i++ {
		dishVo[i] = response.DishListVo{
			Id:          dishes[i].Id,
			Name:        dishes[i].Name,
			CategoryId:  dishes[i].CategoryId,
			Description: dishes[i].Description,
			Price:       dishes[i].Price,
			Image:       dishes[i].Image,
			Status:      dishes[i].Status,
			CreateTime:  dishes[i].CreateTime,
			UpdateTime:  dishes[i].UpdateTime,
			CreateUser:  dishes[i].CreateUser,
			UpdateUser:  dishes[i].UpdateUser,
		}
	}

	// 设置redis缓存

	if dishVoJSON, err := json.Marshal(dishVo); err == nil {
		// 设置缓存，过期时间为0表示不会过期
		if err = global.RedisClient.Set(DishCacheKey+strconv.Itoa(int(categoryId)), dishVoJSON, 0).Err(); err != nil {
			global.Log.Warn("redis 缓存菜品设置失败")
		}
	}

	return dishVo, nil
}

func (d DishServiceImpl) OnOrClose(ctx context.Context, id uint64, status int) error {
	utils.CleanCache(DishCacheKey + "*")
	return d.repo.OnOrClose(ctx, id, status)
}

func (d DishServiceImpl) Update(ctx context.Context, dto request.DishUpdateDTO) error {
	var err error
	price, _ := strconv.ParseFloat(dto.Price, 64)
	dish := model.Dish{
		Id:          dto.Id,
		Name:        dto.Name,
		CategoryId:  dto.CategoryId,
		Price:       price,
		Image:       dto.Image,
		Description: dto.Description,
		Status:      enum.ENABLE,
		Flavors:     dto.Flavors,
	}

	// 事务开启
	transaction := d.repo.Transaction(ctx)
	defer func() {
		if r := recover(); r != nil {
			// 发生 panic 时回滚事务
			transaction.Rollback()
		} else if err != nil {
			// 发生错误时回滚事务
			transaction.Rollback()
		}
	}()

	// 更新菜品信息
	if err = d.repo.Update(transaction, dish); err != nil {
		return err
	}
	// 更新菜品的口味分两步： 1.先删除原有的所有关联数据，2.再插入新的口味数据
	if err = d.dishFlavorRepo.DeleteByDishId(transaction, dish.Id); err != nil {
		return err
	}
	if len(dish.Flavors) != 0 {
		if err = d.dishFlavorRepo.InsertBatch(transaction, dish.Flavors); err != nil {
			return err
		}
	}

	if err = transaction.Commit().Error; err != nil {
		return err // 这里会直接返回错误，defer 中的回滚会执行一次
	}
	utils.CleanCache(DishCacheKey + "*")
	return nil
}

func (d *DishServiceImpl) Delete(ctx context.Context, ids string) error {
	var err error
	// 删除分两步， 删除菜品和删除关联的口味 (ids 为需要删除的菜品id集合，如"1,2,3")
	idList := strings.Split(ids, ",")
	for _, idStr := range idList {
		dishId, _ := strconv.ParseUint(idStr, 10, 64)
		// 开始一个新的事务
		transaction := d.repo.Transaction(ctx)
		defer func() {
			if r := recover(); r != nil {
				transaction.Rollback()
			} else if err != nil {
				transaction.Rollback()
			}
		}()
		// 删除菜品的口味数据
		err = d.dishFlavorRepo.DeleteByDishId(transaction, dishId)
		if err != nil {
			return err
		}
		// 删除菜品
		err = d.repo.Delete(transaction, dishId)
		if err != nil {
			return err
		}
		// 提交事务
		if err = transaction.Commit().Error; err != nil {
			return err
		}
	}
	utils.CleanCache(DishCacheKey + "*")
	return nil
}

func NewDishService(repo repository.DishRepo, dishFlavorRepo repository.DishFlavorRepo) IDishService {
	return &DishServiceImpl{repo: repo, dishFlavorRepo: dishFlavorRepo}
}
