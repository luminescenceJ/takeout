package dao

import (
	"gorm.io/gorm"
	"takeout/internal/model"
	"takeout/repository"
)

type DishFlavorDao struct {
}

func (d DishFlavorDao) InsertBatch(transaction *gorm.DB, flavor []model.DishFlavor) error {
	return transaction.Create(&flavor).Error
}

func (d DishFlavorDao) DeleteByDishId(transaction *gorm.DB, dishId uint64) error {
	return transaction.Where("dish_id = ?", dishId).Delete(&model.DishFlavor{}).Error
}

func (d DishFlavorDao) GetByDishId(transaction *gorm.DB, dishId uint64) ([]model.DishFlavor, error) {
	var flavors []model.DishFlavor
	if err := transaction.Where("dish_id = ?", dishId).Find(&flavors).Error; err != nil {
		return nil, err
	}
	return flavors, nil
}

func (d DishFlavorDao) Update(transaction *gorm.DB, flavor model.DishFlavor) error {
	return transaction.Model(&model.DishFlavor{Id: flavor.Id}).Updates(flavor).Error
}

// NewDishFlavorDao db 是上个事务创建出来的
func NewDishFlavorDao() repository.DishFlavorRepo {
	return &DishFlavorDao{}
}
