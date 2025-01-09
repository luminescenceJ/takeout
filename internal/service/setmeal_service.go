package service

import (
	"context"
	"strconv"
	"takeout/common"
	"takeout/internal/api/request"
	"takeout/internal/api/response"
	"takeout/internal/model"
	"takeout/repository"
)

type ISetMealService interface {
	SaveWithDish(ctx context.Context, dto request.SetMealDTO) error
	PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error)
	OnOrClose(ctx context.Context, id uint64, status int) error
	GetByIdWithDish(ctx context.Context, dishId uint64) (response.SetMealWithDishByIdVo, error)
}

type SetMealServiceImpl struct {
	repo            repository.SetMealRepo
	setMealDishRepo repository.SetMealDishRepo
}

func (s SetMealServiceImpl) SaveWithDish(ctx context.Context, dto request.SetMealDTO) error {
	var (
		err error
	)
	price, _ := strconv.ParseFloat(dto.Price, 64)
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

func NewSetMealService(repo repository.SetMealRepo, setMealDishRepo repository.SetMealDishRepo) ISetMealService {
	return &SetMealServiceImpl{
		repo:            repo,
		setMealDishRepo: setMealDishRepo,
	}
}
