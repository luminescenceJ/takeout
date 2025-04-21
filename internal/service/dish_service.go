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

var bf *utils.RedisBloom

func init() {
	bf = utils.NewRedisBloomFilter(DishCacheKey)
}

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
			global.Log.Error("Transaction panicked", "error", r)
			transaction.Rollback()
		} else if err != nil {
			// 发生错误时回滚事务
			global.Log.Error("Error occurred during transaction", "error", err)
			transaction.Rollback()
		}
	}()

	// 插入菜品
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
		global.Log.Error("Failed to insert dish", "dish", dto, "error", err)
		return err
	}
	global.Log.Info("Successfully inserted dish", "dish", dto)

	// 外键设置
	for i := range dto.Flavors {
		dto.Flavors[i].DishId = dish.Id
	}
	if err = d.dishFlavorRepo.InsertBatch(transaction, dto.Flavors); err != nil {
		global.Log.Error("Failed to insert dish flavors", "dishId", dish.Id, "error", err)
		return err
	}
	global.Log.Info("Successfully inserted dish flavors", "dishId", dish.Id)

	// 提交事务
	if err = transaction.Commit().Error; err != nil {
		global.Log.Error("Failed to commit transaction", "error", err)
		return err
	}

	// 清理缓存
	utils.CleanCache(DishCacheKey + "*")
	// 添加布隆过滤器
	_ = bf.AddDish(strconv.FormatUint(dto.Id, 10))
	global.Log.Info("Cache cleared for dish data")

	return nil
}

func (d DishServiceImpl) PageQuery(ctx context.Context, dto *request.DishPageQueryDTO) (*common.PageResult, error) {
	global.Log.Debug("Received request for dish page query", "params", dto)

	pageResult, err := d.repo.PageQuery(ctx, dto)
	if err != nil {
		global.Log.Error("Failed to query page result", "error", err)
		return nil, err
	}

	global.Log.Info("Successfully fetched page query result", "result", pageResult)
	return pageResult, nil
}

func (d DishServiceImpl) GetByIdWithFlavors(ctx context.Context, id uint64) (response.DishVo, error) {
	var dishVO response.DishVo
	global.Log.Debug("Received request to get dish by ID", "id", id)

	dish, err := d.repo.GetById(ctx, id)
	if err != nil {
		global.Log.Error("Failed to fetch dish by ID", "id", id, "error", err)
		return response.DishVo{}, err
	}

	// 填充数据
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

	global.Log.Info("Successfully fetched dish by ID", "id", id, "dish", dishVO)
	return dishVO, nil
}

func (d DishServiceImpl) List(ctx context.Context, categoryId uint64) ([]response.DishListVo, error) {
	var (
		dishes    []model.Dish
		cacheData string
		dishList  []response.DishListVo
		err       error
	)

	// 查询 Redis 缓存
	cacheData, err = global.RedisClient.Get(DishCacheKey + strconv.Itoa(int(categoryId))).Result()

	if err == nil {
		global.Log.Info("Redis cache hit for dishes", "categoryId", categoryId)
		if err = json.Unmarshal([]byte(cacheData), &dishList); err == nil {
			return dishList, nil
		} else {
			global.Log.Warn("Failed to unmarshal cached data", "categoryId", categoryId, "error", err)
		}
	} else {
		global.Log.Info("Cache miss for dishes", "categoryId", categoryId)
		// 布隆过滤器
		has, _ := bf.HasDish(strconv.FormatUint(categoryId, 10))
		if !has {
			global.Log.Info("Bad Request,No Data in Mysql")
			return make([]response.DishListVo, 0), nil
		}
	}

	// 查询数据库
	dishes, err = d.repo.List(ctx, categoryId)
	if err != nil {
		global.Log.Error("Failed to fetch dishes from DB", "categoryId", categoryId, "error", err)
		return nil, err
	}

	// 转换为响应数据格式
	count := len(dishes)
	dishList = make([]response.DishListVo, count)
	for i := 0; i < count; i++ {
		dishList[i] = response.DishListVo{
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

	// 设置 Redis 缓存
	if dishVoJSON, err := json.Marshal(dishList); err == nil {
		//global.RedisClient.HSet()
		if err = global.RedisClient.Set(DishCacheKey+strconv.Itoa(int(categoryId)), dishVoJSON, 0).Err(); err != nil {
			global.Log.Warn("Failed to set Redis cache", "categoryId", categoryId, "error", err)
		} else {
			global.Log.Info("Successfully cached dishes in Redis", "categoryId", categoryId)
		}
	}

	_ = bf.AddDish(strconv.FormatUint(categoryId, 10))

	return dishList, nil
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
	// 添加布隆过滤器
	_ = bf.AddDish(strconv.FormatUint(dto.Id, 10))
	return nil
}

func (d *DishServiceImpl) Delete(ctx context.Context, ids string) error {
	var err error
	// 删除分两步，删除菜品和删除关联的口味 (ids 为需要删除的菜品id集合，如"1,2,3")
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
			global.Log.Error("Failed to delete dish flavors", "dishId", dishId, "error", err)
			return err
		}

		// 删除菜品
		err = d.repo.Delete(transaction, dishId)
		if err != nil {
			global.Log.Error("Failed to delete dish", "dishId", dishId, "error", err)
			return err
		}

		// 提交事务
		if err = transaction.Commit().Error; err != nil {
			global.Log.Error("Failed to commit transaction", "dishId", dishId, "error", err)
			return err
		}

		// 记录删除成功
		global.Log.Info("Successfully deleted dish", "dishId", dishId)
	}

	// 清理缓存
	utils.CleanCache(DishCacheKey + "*")
	global.Log.Info("Cache cleared after deletion")

	return nil
}

func NewDishService(repo repository.DishRepo, dishFlavorRepo repository.DishFlavorRepo) IDishService {
	return &DishServiceImpl{repo: repo, dishFlavorRepo: dishFlavorRepo}
}
