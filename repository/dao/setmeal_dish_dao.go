package dao

import (
	"gorm.io/gorm"
	"takeout/internal/model"
	"takeout/repository"
)

type SetMealDishDao struct {
}

func (s SetMealDishDao) DeleteBySetMealId(transaction *gorm.DB, SetMealId uint64) error {
	return transaction.Where("setmeal_id=?", SetMealId).Delete(&model.SetMealDish{}).Error
}

func (s SetMealDishDao) InsertBatch(transaction *gorm.DB, setmealDishs []model.SetMealDish) error {
	return transaction.Create(&setmealDishs).Error
}

func (s SetMealDishDao) GetBySetMealId(transaction *gorm.DB, SetMealId uint64) ([]model.SetMealDish, error) {
	var setmealDishs []model.SetMealDish
	err := transaction.Where("setmeal_id = ?", SetMealId).Find(&setmealDishs).Error
	return setmealDishs, err
}

func NewSetMealDishDao() repository.SetMealDishRepo {
	return &SetMealDishDao{}
}
