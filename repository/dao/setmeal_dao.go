package dao

import (
	"context"
	"gorm.io/gorm"
	"takeout/common"
	"takeout/internal/api/request"
	"takeout/internal/api/response"
	"takeout/internal/model"
	"takeout/repository"
)

type SetMealDao struct {
	db *gorm.DB
}

func (s SetMealDao) Update(transaction *gorm.DB, meal *model.SetMeal) error {
	return transaction.Updates(&meal).Error
}

func (s SetMealDao) Delete(transaction *gorm.DB, id uint64) error {
	return transaction.Where("id = ?", id).Delete(&model.SetMeal{}).Error
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
	var (
		res    common.PageResult
		record []response.SetMealPageQueryVo
		err    error
	)
	query := s.db.WithContext(ctx).Model(&model.SetMeal{})

	if dto.CategoryId != 0 {
		query = query.Where("setmeal.category_id = ?", dto.CategoryId)
	}
	if dto.Name != "" {
		query = query.Where("setmeal.name LIKE ?", "%"+dto.Name+"%")
	}
	if dto.Status != "" {
		query = query.Where("setmeal.status = ?", dto.Status)
	}
	if err = query.Count(&res.Total).Error; err != nil {
		return nil, err
	}

	err = query.Scopes(res.Paginate(&dto.Page, &dto.PageSize)).
		Select("setmeal.* ,c.name as category_name").
		Joins("LEFT JOIN category c on setmeal.category_id=c.id").
		Order("create_time desc").
		Scan(&record).
		Error
	res.Records = record

	return &res, err
}

func (s SetMealDao) SetStatus(ctx context.Context, id uint64, status int) error {
	return s.db.WithContext(ctx).Model(&model.SetMeal{}).Where("id=?", id).Update("status", status).Error
}

func (s SetMealDao) GetByIdWithDish(db *gorm.DB, dishId uint64) (model.SetMeal, error) {

	var (
		setMeal model.SetMeal
		err     error
	)
	if err = db.Where("id=?", dishId).Find(&setMeal).Error; err != nil {
		return model.SetMeal{}, err
	}
	return setMeal, nil

}

func NewSetMealDao(db *gorm.DB) repository.SetMealRepo {
	return &SetMealDao{db: db}
}
