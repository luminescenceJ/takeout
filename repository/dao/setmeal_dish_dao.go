package dao

import (
	"takeout/global/tx"
	"takeout/internal/model"
	"takeout/repository"
)

type SetMealDishDao struct {
}

func (s SetMealDishDao) InsertBatch(db tx.Transaction, setmealDishs []model.SetMealDish) error {
	//TODO implement me
	panic("implement me")
}

func (s SetMealDishDao) GetBySetMealId(db tx.Transaction, SetMealId uint64) ([]model.SetMealDish, error) {
	//TODO implement me
	panic("implement me")
}

func NewSetMealDishDao() repository.SetMealDishRepo {
	return &SetMealDishDao{}
}
