package service

import (
	"context"
	"strconv"
	"takeout/common"
	"takeout/common/enum"
	"takeout/internal/api/request"
	"takeout/internal/model"
	"takeout/repository"
)

type ICategoryService interface {
	AddCategory(ctx context.Context, dto request.CategoryDTO) error
	PageQuery(ctx context.Context, dto request.CategoryPageQueryDTO) (*common.PageResult, error)
	List(ctx context.Context, cate int) ([]model.Category, error)
	DeleteById(ctx context.Context, id uint64) error
	Update(ctx context.Context, dto request.CategoryDTO) error
	SetStatus(ctx context.Context, id uint64, status int) error
}

type CategoryImpl struct {
	repo repository.CategoryRepo
}

func (c *CategoryImpl) AddCategory(ctx context.Context, dto request.CategoryDTO) error {
	typeInStr, _ := strconv.Atoi(dto.Type)
	sortInStr, _ := strconv.Atoi(dto.Sort)
	return c.repo.Insert(ctx, model.Category{
		Type:   typeInStr,
		Name:   dto.Name,
		Sort:   sortInStr,
		Status: enum.DISABLE,
	})
}

func (c *CategoryImpl) PageQuery(ctx context.Context, dto request.CategoryPageQueryDTO) (*common.PageResult, error) {
	return c.repo.PageQuery(ctx, dto)
}

func (c *CategoryImpl) List(ctx context.Context, cate int) ([]model.Category, error) {

	return c.repo.List(ctx, cate)
}

func (c *CategoryImpl) DeleteById(ctx context.Context, id uint64) error {
	return c.repo.DeleteById(ctx, id)
}

func (c *CategoryImpl) Update(ctx context.Context, dto request.CategoryDTO) error {
	sort, _ := strconv.Atoi(dto.Sort)
	type_, _ := strconv.Atoi(dto.Type)
	return c.repo.Update(ctx, model.Category{
		Id:   dto.Id,
		Name: dto.Name,
		Sort: sort,
		Type: type_,
	})
}

func (c *CategoryImpl) SetStatus(ctx context.Context, id uint64, status int) error {
	return c.repo.SetStatus(ctx, model.Category{
		Id:     id,
		Status: status,
	})
}

func NewCategoryService(repo repository.CategoryRepo) ICategoryService {
	return &CategoryImpl{repo: repo}
}
