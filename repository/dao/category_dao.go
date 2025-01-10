package dao

import (
	"context"
	"gorm.io/gorm"
	"takeout/common"
	"takeout/internal/api/admin/request"
	"takeout/internal/model"
)

type CategoryDao struct {
	db *gorm.DB
}

func (c *CategoryDao) SetStatus(ctx context.Context, category model.Category) error {
	return c.db.WithContext(ctx).Model(&category).Update("status", category.Status).Error
}

func (c *CategoryDao) Update(ctx context.Context, category model.Category) error {
	return c.db.WithContext(ctx).Model(&category).Updates(&category).Error
}

func (c *CategoryDao) DeleteById(ctx context.Context, id uint64) error {
	return c.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

func (c *CategoryDao) List(ctx context.Context, cate int) ([]model.Category, error) {
	var res []model.Category
	err := c.db.WithContext(ctx).
		Where("type = ?", cate).
		Order("sort asc").
		Order("create_time desc").
		Find(&res).
		Error
	return res, err
}

func (c *CategoryDao) PageQuery(ctx context.Context, dto request.CategoryPageQueryDTO) (*common.PageResult, error) {
	var (
		res     common.PageResult
		records []model.Category
		err     error
	)
	query := c.db.WithContext(ctx).Model(model.Category{})
	if dto.Name != "" {
		query = query.Where("name like ?", "%"+dto.Name+"%")
	}
	if dto.Cate != 0 {
		query = query.Where("type = ?", dto.Cate)
	}

	if err = query.Count(&res.Total).Error; err != nil {
		return nil, err
	}

	err = query.Scopes(res.Paginate(&dto.Page, &dto.PageSize)).
		Order("create_time desc").
		Find(&records).
		Error
	//err = query.Offset((dto.Page - 1) * dto.PageSize).Limit(dto.PageSize).Find(&records).Error

	res.Records = records
	return &res, err
}

func (c *CategoryDao) Insert(ctx context.Context, category model.Category) error {
	return c.db.WithContext(ctx).Create(&category).Error
}

func NewCategoryDao(db *gorm.DB) *CategoryDao {
	return &CategoryDao{db: db}
}
