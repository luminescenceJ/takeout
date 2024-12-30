package dao

import (
	"context"
	"gorm.io/gorm"
	"takeout/common"
	"takeout/internal/api/request"
	"takeout/internal/model"
)

type CategoryDao struct {
	db *gorm.DB
}

func (c *CategoryDao) SetStatus(ctx context.Context, category model.Category) error {
	return nil
}

func (c *CategoryDao) Update(ctx context.Context, category model.Category) error {
	return nil
}

func (c *CategoryDao) DeleteById(ctx context.Context, id uint64) error {
	return nil
}

func (c *CategoryDao) List(ctx context.Context, cate int) ([]model.Category, error) {

	return nil, nil
}

func (c *CategoryDao) PageQuery(ctx context.Context, dto request.CategoryPageQueryDTO) (*common.PageResult, error) {

	return nil, nil
}

func (c *CategoryDao) Insert(ctx context.Context, category model.Category) error {
	return c.db.WithContext(ctx).Create(&category).Error
}

func NewCategoryDao(db *gorm.DB) *CategoryDao {
	return &CategoryDao{db: db}
}
