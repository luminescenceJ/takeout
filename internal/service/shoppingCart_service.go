package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/deepcopier"
	"takeout/common/enum"
	"takeout/internal/api/user/request"
	"takeout/internal/model"
	"takeout/repository"
)

type IShoppingCartService interface {
	AddShoppingCart(ctx *gin.Context, dto request.ShoppingCartDTO) error
	QueryShoppingCart(ctx *gin.Context) ([]model.ShoppingCart, error)
	CleanShoppingCart(ctx *gin.Context) error //清空
	SubShoppingCart(ctx *gin.Context, ShoppingCartDTO request.ShoppingCartDTO) error
}

type ShoppingCartService struct {
	repo repository.ShoppingCartRepo
}

func (s ShoppingCartService) AddShoppingCart(ctx *gin.Context, dto request.ShoppingCartDTO) error {
	var (
		userId       uint64
		shoppingCart model.ShoppingCart
		err          error
	)
	if err = deepcopier.Copy(dto).To(&shoppingCart); err != nil {
		return err
	}
	if CurrentId, ok := ctx.Get(enum.CurrentId); ok {
		userId = CurrentId.(uint64)
	} else {
		return errors.New("用户不存在")
	}
	shoppingCart.UserId = int(userId)
	if err = s.repo.Add(ctx, shoppingCart); err != nil {
		return err
	}
	return nil
}

func (s ShoppingCartService) QueryShoppingCart(ctx *gin.Context) ([]model.ShoppingCart, error) {
	var userId uint64
	if CurrentId, ok := ctx.Get(enum.CurrentId); ok {
		userId = CurrentId.(uint64)
	}
	return s.repo.List(ctx, int(userId))
}

func (s ShoppingCartService) CleanShoppingCart(ctx *gin.Context) error {
	var userId uint64
	if CurrentId, ok := ctx.Get(enum.CurrentId); ok {
		userId = CurrentId.(uint64)
	}
	return s.repo.Clean(ctx, int(userId))
}

func (s ShoppingCartService) SubShoppingCart(ctx *gin.Context, ShoppingCartDTO request.ShoppingCartDTO) error {
	var userId uint64
	if CurrentId, ok := ctx.Get(enum.CurrentId); ok {
		userId = CurrentId.(uint64)
	}
	return s.repo.Subtract(ctx, ShoppingCartDTO, int(userId))
}

func NewShoppingCartService(repo repository.ShoppingCartRepo) IShoppingCartService {
	return &ShoppingCartService{repo: repo}
}
