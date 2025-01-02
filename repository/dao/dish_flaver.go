package dao

import (
	"gorm.io/gorm"
	"takeout/internal/model"
	"takeout/repository"
)

type DishFlavorDao struct {
}

func (d DishFlavorDao) InsertBatch(db *gorm.DB, flavor []model.DishFlavor) error {
	//TODO implement me
	panic("implement me")
}

func (d DishFlavorDao) DeleteByDishId(db *gorm.DB, dishId uint64) error {
	//TODO implement me
	panic("implement me")
}

func (d DishFlavorDao) GetByDishId(db *gorm.DB, dishId uint64) ([]model.DishFlavor, error) {
	//TODO implement me
	panic("implement me")
}

func (d DishFlavorDao) Update(db *gorm.DB, flavor model.DishFlavor) error {
	//TODO implement me
	panic("implement me")
}

// NewDishFlavorDao db 是上个事务创建出来的
func NewDishFlavorDao() repository.DishFlavorRepo {
	return &DishFlavorDao{}
}
