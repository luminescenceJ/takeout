package repository

import (
	"context"
	"takeout/internal/api/user/request"
	"takeout/internal/model"
)

type ShoppingCartRepo interface {
	Add(ctx context.Context, shoppingCart model.ShoppingCart) error
	List(ctx context.Context, userId int) ([]model.ShoppingCart, error)
	Clean(ctx context.Context, userId int) error //清空
	Subtract(ctx context.Context, ShoppingCartDTO request.ShoppingCartDTO, userId int) error
	InsertBatchShoppingCart(ctx context.Context, shoppingCartList []model.ShoppingCart) error
}
