package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/user/request"
	"takeout/internal/model"
	"takeout/internal/service"
)

type ShoppingCartController struct {
	service service.IShoppingCartService
}

func NewShoppingController(service service.IShoppingCartService) *ShoppingCartController {
	return &ShoppingCartController{service: service}
}

// AddShoppingCart @AddShoppingCart 添加购物车
// @Tags ShoppingCart
// @Security JWTAuth
// @Produce json
// @Param model body model.ShoppingCartDTO true "菜品属性"
// @Success 200 {object} common.Result{} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/shoppingCart/add [post]
func (s *ShoppingCartController) AddShoppingCart(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		dto  request.ShoppingCartDTO
		err  error
	)

	if err = ctx.ShouldBindJSON(&dto); err != nil {
		global.Log.Debug("C端-购物车接口 AddShoppingCart param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = s.service.AddShoppingCart(ctx, dto); err != nil {
		code = e.ERROR
		global.Log.Debug("AddShoppingCart error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})

}

// QueryShoppingCart @QueryShoppingCart 查看购物车
// @Tags ShoppingCart
// @Security JWTAuth
// @Produce json
// @Success 200 {object} common.Result{Data=[]model.ShoppingCart} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/shoppingCart/list [get]
func (s *ShoppingCartController) QueryShoppingCart(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		data []model.ShoppingCart
		err  error
	)

	if data, err = s.service.QueryShoppingCart(ctx); err != nil {
		code = e.ERROR
		global.Log.Debug("AddShoppingCart error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: data,
		Msg:  e.GetMsg(code),
	})
}

// CleanShoppingCart @CleanShoppingCart 添加购物车
// @Tags ShoppingCart
// @Security JWTAuth
// @Produce json
// @Success 200 {object} common.Result{} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/shoppingCart/clean [delete]
func (s *ShoppingCartController) CleanShoppingCart(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		dto  request.ShoppingCartDTO
		err  error
	)

	if err = ctx.ShouldBindJSON(&dto); err != nil {
		global.Log.Debug("C端-购物车接口 AddShoppingCart param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = s.service.CleanShoppingCart(ctx); err != nil {
		code = e.ERROR
		global.Log.Debug("AddShoppingCart error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// SubShoppingCart @SubShoppingCart 删除购物车中一个商品
// @Tags ShoppingCart
// @Security JWTAuth
// @Produce json
// @Param dto body request.ShoppingCartDTO true "菜品属性"
// @Success 200 {object} common.Result{Data=response.DishListVo} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/shoppingCart/sub [post]
func (s *ShoppingCartController) SubShoppingCart(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		dto  request.ShoppingCartDTO
		err  error
	)

	if err = ctx.ShouldBindJSON(&dto); err != nil {
		global.Log.Debug("C端-购物车接口 AddShoppingCart param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = s.service.SubShoppingCart(ctx, dto); err != nil {
		code = e.ERROR
		global.Log.Debug("AddShoppingCart error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}
