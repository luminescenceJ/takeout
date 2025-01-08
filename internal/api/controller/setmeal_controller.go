package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/request"
	"takeout/internal/service"
)

type SetMealController struct {
	service service.ISetMealService
}

func NewSetMealController(service service.ISetMealService) *SetMealController {
	return &SetMealController{service: service}
}

// SaveWithDish @SaveWithDish 保存套餐和菜品信息
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param dto body request.SetMealDTO true
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/setmeal [post]
func (sc *SetMealController) SaveWithDish(ctx *gin.Context) {
	var (
		dto  request.SetMealDTO
		code = e.SUCCESS
		err  error
	)

	if err = ctx.ShouldBindJSON(&dto); err != nil {
		global.Log.Debug("SaveWithDish保存套餐和菜品信息 结构体解析失败！")
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = sc.service.SaveWithDish(ctx, dto); err != nil {
		code = e.ERROR
		global.Log.Warn("AddDish failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// PageQuery 套餐分页查询
func (sc *SetMealController) PageQuery(ctx *gin.Context) {
	code := e.SUCCESS

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// OnOrClose 套餐启用禁用
func (sc *SetMealController) OnOrClose(ctx *gin.Context) {
	code := e.SUCCESS

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// GetByIdWithDish 根据套餐id获取套餐和关联菜品信息
func (sc *SetMealController) GetByIdWithDish(ctx *gin.Context) {
	code := e.SUCCESS

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}
