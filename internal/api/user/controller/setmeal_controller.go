package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	userResponse "takeout/internal/api/user/response"
	"takeout/internal/model"
	"takeout/internal/service"
)

type SetMealController struct {
	service service.ISetMealService
}

func NewSetMealController(service service.ISetMealService) *SetMealController {
	return &SetMealController{service: service}
}

// GetByIdWithDish @GetByIdWithDish 根据分类id查询套餐
// @Tags UserSetmeal
// @Security JWTAuth
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/setmeal/dish/{id} [get]
func (sc *SetMealController) GetDishByCategoryId(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		id      string
		setmeal []userResponse.DishItemVO
		err     error
	)
	id = ctx.Param("id")
	uint64Id, _ := strconv.ParseUint(id, 10, 64)
	if setmeal, err = sc.service.GetDishBySetmealId(ctx, uint64Id); err != nil {
		code = e.ERROR
		global.Log.Warn("GetDishByCategoryId failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Data: nil,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: setmeal,
		Msg:  e.GetMsg(code),
	})
}

// List @List 根据套餐id查询包含的菜品
// @Tags UserSetmeal
// @Security JWTAuth
// @Produce json
// @Param categoryId query string true "分类id"
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/setmeal/list [get]
func (sc *SetMealController) List(ctx *gin.Context) {
	var (
		code       = e.SUCCESS
		categoryId string
		setmeals   []model.SetMeal
		err        error
	)
	categoryId = ctx.Query("categoryId")
	if setmeals, err = sc.service.List(ctx, categoryId); err != nil {
		code = e.ERROR
		global.Log.Warn("List failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Data: nil,
			Msg:  e.GetMsg(code),
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: setmeals,
		Msg:  e.GetMsg(code),
	})
}
