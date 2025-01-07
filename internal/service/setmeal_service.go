package service

import (
	"context"
	"takeout/common"
	"takeout/internal/api/request"
	"takeout/internal/api/response"
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
	//TODO implement me
	panic("implement me")
}

func (s SetMealServiceImpl) PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s SetMealServiceImpl) OnOrClose(ctx context.Context, id uint64, status int) error {
	//TODO implement me
	panic("implement me")
}

func (s SetMealServiceImpl) GetByIdWithDish(ctx context.Context, dishId uint64) (response.SetMealWithDishByIdVo, error) {
	//TODO implement me
	panic("implement me")
}

func NewSetMealService(repo repository.SetMealRepo, setMealDishRepo repository.SetMealDishRepo) ISetMealService {
	return &SetMealServiceImpl{
		repo:            repo,
		setMealDishRepo: setMealDishRepo,
	}
}
