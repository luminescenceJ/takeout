package repository

import (
	"gorm.io/gorm"
	"takeout/internal/model"
)

type SetMealDishRepo interface {
	InsertBatch(db *gorm.DB, setmealDishs []model.SetMealDish) error
	GetBySetMealId(db *gorm.DB, SetMealId uint64) ([]model.SetMealDish, error)
	DeleteBySetMealId(db *gorm.DB, SetMealId uint64) error
}
