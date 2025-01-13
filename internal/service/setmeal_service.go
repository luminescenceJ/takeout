package service

import (
	"context"
	"strconv"
	"strings"
	"takeout/common"
	"takeout/internal/api/admin/request"
	"takeout/internal/api/admin/response"
	userResponse "takeout/internal/api/user/response"
	"takeout/internal/model"
	"takeout/repository"
)

type ISetMealService interface {
	SaveWithDish(ctx context.Context, dto request.SetMealDTO) error
	PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error)
	OnOrClose(ctx context.Context, id uint64, status int) error
	GetByIdWithDish(ctx context.Context, dishId uint64) (response.SetMealWithDishByIdVo, error)
	Update(ctx context.Context, dto request.SetMealDTO) error
	DeleteBatch(ctx context.Context, ids string) error
	GetDishBySetmealId(ctx context.Context, setmealId uint64) ([]userResponse.DishItemVO, error)
	List(ctx context.Context, categoryId string) ([]model.SetMeal, error)
}

type SetMealServiceImpl struct {
	repo            repository.SetMealRepo
	setMealDishRepo repository.SetMealDishRepo
}

func (s SetMealServiceImpl) Update(ctx context.Context, dto request.SetMealDTO) error {
	var (
		err   error
		price = float64(dto.Price)
		meal  = &model.SetMeal{
			Id:          dto.Id,
			CategoryId:  dto.CategoryId,
			Name:        dto.Name,
			Price:       price,
			Status:      dto.Status,
			Description: dto.Description,
			Image:       dto.Image,
		}
	)

	// 开启事务 先更新套餐setmeal表， 再更新其对应的菜品setmeal_dish表
	transaction := s.repo.Transaction(ctx)
	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		} else if err != nil {
			transaction.Rollback()
		}
	}()
	// 更新套餐
	if err = s.repo.Update(transaction, meal); err != nil {
		return err
	}
	// 删除原来绑定的菜品
	if err = s.setMealDishRepo.DeleteBySetMealId(transaction, dto.Id); err != nil {
		return err
	}
	// 绑定新的菜品
	for i := range dto.SetMealDishs {
		dto.SetMealDishs[i].SetmealId = dto.Id
	}
	if err = s.setMealDishRepo.InsertBatch(transaction, dto.SetMealDishs); err != nil {
		return err
	}
	// 提交事务
	if err = transaction.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s SetMealServiceImpl) DeleteBatch(ctx context.Context, ids string) error {
	var (
		idSet    = strings.Split(ids, ",")
		deleteId uint64
		err      error
	)
	// 对于每个需要删除的套餐单独开启事务
	// 先删除套餐再删除套餐id对应的菜品
	// 最后提交事务，如果出现错误直接回滚该事务
	for _, id := range idSet {
		transaction := s.repo.Transaction(ctx)
		defer func() {
			if r := recover(); r != nil {
				transaction.Rollback()
			} else if err != nil {
				transaction.Rollback()
			}
		}()
		if deleteId, err = strconv.ParseUint(id, 10, 64); err != nil {
			return err
		}
		if err = s.repo.Delete(transaction, deleteId); err != nil {
			return err
		}
		if err = s.setMealDishRepo.DeleteBySetMealId(transaction, deleteId); err != nil {
			return err
		}
		if err = transaction.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}

func (s SetMealServiceImpl) SaveWithDish(ctx context.Context, dto request.SetMealDTO) error {
	var (
		err   error
		price = float64(dto.Price)
	)
	setmeal := &model.SetMeal{
		CategoryId:  dto.CategoryId,
		Name:        dto.Name,
		Price:       price,
		Status:      dto.Status,
		Description: dto.Description,
		Image:       dto.Image,
	}

	transaction := s.repo.Transaction(ctx)
	defer func() {
		if r := recover(); r != nil {
			// 发生 panic 时回滚事务
			transaction.Rollback()
		} else if err != nil {
			// 发生错误时回滚事务
			transaction.Rollback()
		}
	}()

	if err = s.repo.Insert(transaction, setmeal); err != nil {
		return err
	}

	for i := range dto.SetMealDishs {
		dto.SetMealDishs[i].SetmealId = setmeal.Id
	}
	if err = s.setMealDishRepo.InsertBatch(transaction, dto.SetMealDishs); err != nil {
		return err
	}
	if err = transaction.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s SetMealServiceImpl) PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error) {
	return s.repo.PageQuery(ctx, dto)
}

func (s SetMealServiceImpl) OnOrClose(ctx context.Context, id uint64, status int) error {
	return s.repo.SetStatus(ctx, id, status)
}

func (s SetMealServiceImpl) GetByIdWithDish(ctx context.Context, mealId uint64) (response.SetMealWithDishByIdVo, error) {
	var (
		err      error
		res      response.SetMealWithDishByIdVo
		setmeal  model.SetMeal
		dishList []model.SetMealDish
	)
	// 为了保持查询结果的一致性，开启事务
	transaction := s.repo.Transaction(ctx)
	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		} else if err != nil {
			transaction.Rollback()
		}
	}()
	// 两次查询 先查询套餐 再查询关联的菜品
	if setmeal, err = s.repo.GetByIdWithDish(transaction, mealId); err != nil {
		return res, err
	}
	if dishList, err = s.setMealDishRepo.GetBySetMealId(transaction, mealId); err != nil {
		return res, err
	}
	if err = transaction.Commit().Error; err != nil {
		return res, err
	}
	res = response.SetMealWithDishByIdVo{
		Id:            setmeal.Id,
		CategoryId:    setmeal.CategoryId,
		Description:   setmeal.Description,
		Image:         setmeal.Image,
		Name:          setmeal.Name,
		Price:         setmeal.Price,
		Status:        setmeal.Status,
		UpdateTime:    setmeal.UpdateTime,
		CategoryName:  setmeal.Name,
		SetmealDishes: dishList,
	}
	return res, nil
}

func (s SetMealServiceImpl) GetDishBySetmealId(ctx context.Context, setmealId uint64) ([]userResponse.DishItemVO, error) {
	// 根据套餐id查询包含的菜品
	return s.repo.GetDishBySetmealId(ctx, setmealId)
}

func (s SetMealServiceImpl) List(ctx context.Context, categoryId string) ([]model.SetMeal, error) {
	// 根据分类id查询套餐
	id, _ := strconv.ParseUint(categoryId, 10, 64)
	return s.repo.GetSetmealByCategoryId(ctx, id)
}

func NewSetMealService(repo repository.SetMealRepo, setMealDishRepo repository.SetMealDishRepo) ISetMealService {
	return &SetMealServiceImpl{
		repo:            repo,
		setMealDishRepo: setMealDishRepo,
	}
}
