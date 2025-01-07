package repository

import (
	"takeout/global/tx"
	"takeout/internal/model"
)

type SetMealDishRepo interface {
	InsertBatch(db tx.Transaction, setmealDishs []model.SetMealDish) error
	GetBySetMealId(db tx.Transaction, SetMealId uint64) ([]model.SetMealDish, error)
}
