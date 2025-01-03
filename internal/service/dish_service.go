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
	price, _ := strconv.ParseFloat(dto.Price, 64)
	transaction := d.repo.Transaction(ctx)
	defer func() {
		if r := recover(); r != nil {
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
		Status:      dto.Status,
		Flavors:     dto.Flavors,
	}
	if err := d.repo.Insert(transaction, &dish); err != nil {
		return err
	}
	for i := range dto.Flavors {
		dto.Flavors[i].DishId = dish.Id
	}
	if err := d.dishFlavorRepo.InsertBatch(transaction, dto.Flavors); err != nil {
		return err
	}
	return transaction.Commit().Error
}

func (d DishServiceImpl) PageQuery(ctx context.Context, dto *request.DishPageQueryDTO) (*common.PageResult, error) {
	//TODO implement me
	panic("implement me")
}

func (d DishServiceImpl) GetByIdWithFlavors(ctx context.Context, id uint64) (response.DishVo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DishServiceImpl) List(ctx context.Context, categoryId uint64) ([]response.DishListVo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DishServiceImpl) OnOrClose(ctx context.Context, id uint64, status int) error {
	//TODO implement me
	panic("implement me")
}

func (d DishServiceImpl) Update(ctx context.Context, dto request.DishUpdateDTO) error {
	//TODO implement me
	panic("implement me")
}

func (d DishServiceImpl) Delete(ctx context.Context, ids string) error {
	//TODO implement me
	panic("implement me")
}

func NewDishService(repo repository.DishRepo, dishFlavorRepo repository.DishFlavorRepo) IDishService {
	return &DishServiceImpl{repo: repo, dishFlavorRepo: dishFlavorRepo}
}
