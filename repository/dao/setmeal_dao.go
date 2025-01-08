package dao

import (
	"context"
	"gorm.io/gorm"
	"takeout/common"
	"takeout/internal/api/request"
	"takeout/internal/model"
	"takeout/repository"
)

type SetMealDao struct {
	db *gorm.DB
}

func (s SetMealDao) Transaction(ctx context.Context) *gorm.DB {
	var count int64
	_ = s.db.WithContext(ctx).Model(&model.SetMeal{}).Count(&count).Error
	return s.db.WithContext(ctx).Begin()
}

func (s SetMealDao) Insert(transaction *gorm.DB, meal *model.SetMeal) error {
	return transaction.Create(&meal).Error
}

func (s SetMealDao) PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s SetMealDao) SetStatus(ctx context.Context, id uint64, status int) error {
	//TODO implement me
	panic("implement me")
}

func (s SetMealDao) GetByIdWithDish(transaction *gorm.DB, dishId uint64) (model.SetMeal, error) {
	//TODO implement me
	panic("implement me")
}

func NewSetMealDao(db *gorm.DB) repository.SetMealRepo {
	return &SetMealDao{db: db}
}
