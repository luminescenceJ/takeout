package repository

import (
	"context"
	"takeout/common"
	"takeout/global/tx"
	"takeout/internal/api/request"
	"takeout/internal/model"
)

type SetMealRepo interface {
	Transaction(ctx context.Context) tx.Transaction
	Insert(db tx.Transaction, meal *model.SetMeal) error
	PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error)
	SetStatus(ctx context.Context, id uint64, status int) error
	GetByIdWithDish(db tx.Transaction, dishId uint64) (model.SetMeal, error)
}
