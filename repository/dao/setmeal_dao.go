package dao

import (
	"context"
	"gorm.io/gorm"
	"takeout/common"
	"takeout/global/tx"
	"takeout/internal/api/request"
	"takeout/internal/model"
	"takeout/repository"
)

type SetMealDao struct {
	db *gorm.DB
}

func (s SetMealDao) Transaction(ctx context.Context) tx.Transaction {
	//TODO implement me
	panic("implement me")
}

func (s SetMealDao) Insert(db tx.Transaction, meal *model.SetMeal) error {
	//TODO implement me
	panic("implement me")
}

func (s SetMealDao) PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s SetMealDao) SetStatus(ctx context.Context, id uint64, status int) error {
	//TODO implement me
	panic("implement me")
}

func (s SetMealDao) GetByIdWithDish(db tx.Transaction, dishId uint64) (model.SetMeal, error) {
	//TODO implement me
	panic("implement me")
}

func NewSetMealDao(db *gorm.DB) repository.SetMealRepo {
	return &SetMealDao{db: db}
}
