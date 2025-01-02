package dao

import (
	"context"
	"gorm.io/gorm"
	"takeout/common"
	"takeout/internal/api/request"
	"takeout/internal/model"
	"takeout/repository"
)

type DishDao struct {
	db *gorm.DB
}

func (d DishDao) Transaction(ctx context.Context) *gorm.DB {
	//TODO implement me
	panic("implement me")
}

func (d DishDao) Insert(db *gorm.DB, dish *model.Dish) error {
	//TODO implement me
	panic("implement me")
}

func (d DishDao) PageQuery(ctx context.Context, dto *request.DishPageQueryDTO) (*common.PageResult, error) {
	//TODO implement me
	panic("implement me")
}

func (d DishDao) GetById(ctx context.Context, id uint64) (*model.Dish, error) {
	//TODO implement me
	panic("implement me")
}

func (d DishDao) List(ctx context.Context, categoryId uint64) ([]model.Dish, error) {
	//TODO implement me
	panic("implement me")
}

func (d DishDao) OnOrClose(ctx context.Context, id uint64, status int) error {
	//TODO implement me
	panic("implement me")
}

func (d DishDao) Update(db *gorm.DB, dish model.Dish) error {
	//TODO implement me
	panic("implement me")
}

func (d DishDao) Delete(db *gorm.DB, id uint64) error {
	//TODO implement me
	panic("implement me")
}

func NewDishRepo(db *gorm.DB) repository.DishRepo {
	return &DishDao{db: db}
}
