package service

import (
	"context"
	"errors"
	"strconv"
	"takeout/common"
	"takeout/common/enum"
	"takeout/internal/api/admin/request"
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
	typeInStr, err := strconv.Atoi(dto.Type)
	if err != nil {
		return errors.New("invalid category type")
	}
	sortInStr, err := strconv.Atoi(dto.Sort)
	if err != nil {
		return errors.New("invalid category sort")
	}

	if dto.Name == "" {
		return errors.New("category name cannot be empty")
	}

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
	if id == 0 {
		return errors.New("invalid category ID")
	}
	return c.repo.DeleteById(ctx, id)
}

func (c *CategoryImpl) Update(ctx context.Context, dto request.CategoryDTO) error {
	if dto.Name == "" {
		return errors.New("category name cannot be empty")
	}
	sort, err := strconv.Atoi(dto.Sort)
	if err != nil {
		return errors.New("invalid category sort")
	}
	type_, err := strconv.Atoi(dto.Type)
	if err != nil {
		return errors.New("invalid category type")
	}
	return c.repo.Update(ctx, model.Category{
		Id:   dto.Id,
		Name: dto.Name,
		Sort: sort,
		Type: type_,
	})
}

func (c *CategoryImpl) SetStatus(ctx context.Context, id uint64, status int) error {
	if status != 0 && status != 1 {
		return errors.New("invalid status, it must be 0 or 1")
	}
	return c.repo.SetStatus(ctx, model.Category{
		Id:     id,
		Status: status,
	})
}

func NewCategoryService(repo repository.CategoryRepo) ICategoryService {
	return &CategoryImpl{repo: repo}
}

//package service
//
//import (
//	"context"
//	"strconv"
//	"takeout/common"
//	"takeout/common/enum"
//	"takeout/internal/api/admin/request"
//	"takeout/internal/model"
//	"takeout/repository"
//)
//
//type ICategoryService interface {
//	AddCategory(ctx context.Context, dto request.CategoryDTO) error
//	PageQuery(ctx context.Context, dto request.CategoryPageQueryDTO) (*common.PageResult, error)
//	List(ctx context.Context, cate int) ([]model.Category, error)
//	DeleteById(ctx context.Context, id uint64) error
//	Update(ctx context.Context, dto request.CategoryDTO) error
//	SetStatus(ctx context.Context, id uint64, status int) error
//}
//
//type CategoryImpl struct {
//	repo repository.CategoryRepo
//}
//
//func (c *CategoryImpl) AddCategory(ctx context.Context, dto request.CategoryDTO) error {
//	typeInStr, _ := strconv.Atoi(dto.Type)
//	sortInStr, _ := strconv.Atoi(dto.Sort)
//	return c.repo.Insert(ctx, model.Category{
//		Type:   typeInStr,
//		Name:   dto.Name,
//		Sort:   sortInStr,
//		Status: enum.DISABLE,
//	})
//}
//
//func (c *CategoryImpl) PageQuery(ctx context.Context, dto request.CategoryPageQueryDTO) (*common.PageResult, error) {
//	return c.repo.PageQuery(ctx, dto)
//}
//
//func (c *CategoryImpl) List(ctx context.Context, cate int) ([]model.Category, error) {
//
//	return c.repo.List(ctx, cate)
//}
//
//func (c *CategoryImpl) DeleteById(ctx context.Context, id uint64) error {
//	return c.repo.DeleteById(ctx, id)
//}
//
//func (c *CategoryImpl) Update(ctx context.Context, dto request.CategoryDTO) error {
//	sort, _ := strconv.Atoi(dto.Sort)
//	type_, _ := strconv.Atoi(dto.Type)
//	return c.repo.Update(ctx, model.Category{
//		Id:   dto.Id,
//		Name: dto.Name,
//		Sort: sort,
//		Type: type_,
//	})
//}
//
//func (c *CategoryImpl) SetStatus(ctx context.Context, id uint64, status int) error {
//	return c.repo.SetStatus(ctx, model.Category{
//		Id:     id,
//		Status: status,
//	})
//}
//
//func NewCategoryService(repo repository.CategoryRepo) ICategoryService {
//	return &CategoryImpl{repo: repo}
//}
