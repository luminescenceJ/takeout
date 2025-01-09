package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/request"
	"takeout/internal/api/response"
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

// PageQuery @PageQuery 套餐分页查询
// @Tags setmeal
// @Security JWTAuth
// @Produce json
// @Param dto body request.SetMealPageQueryDTO true
// @Success 200 {object} common.Result{data=response.SetMealPageQueryVo} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/setmeal/page [post]
func (sc *SetMealController) PageQuery(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		dto  request.SetMealPageQueryDTO
		res  *common.PageResult
		err  error
	)

	if err = ctx.ShouldBindQuery(&dto); err != nil {
		global.Log.Debug("PageQuery 套餐分页查询 结构体解析失败！:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}
	if res, err = sc.service.PageQuery(ctx, dto); err != nil {
		code = e.ERROR
		global.Log.Warn("PageQuery failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: res,
		Msg:  e.GetMsg(code),
	})
}

// OnOrClose @OnOrClose 套餐启用禁用
// @Tags setmeal
// @Security JWTAuth
// @Produce json
// @Param id query string true
// @Param status path string true
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/setmeal/status/{status} [post]
func (sc *SetMealController) OnOrClose(ctx *gin.Context) {
	var (
		code   = e.SUCCESS
		id     string
		status string
		err    error
	)
	id = ctx.Query("id")
	status = ctx.Param("status")
	if id == "" || status == "" {
		global.Log.Debug("OnOrClose 套餐启用禁用 解析失败！")
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	uint64Id, _ := strconv.ParseUint(id, 10, 64)
	intStatus, _ := strconv.Atoi(status)
	if err = sc.service.OnOrClose(ctx, uint64Id, intStatus); err != nil {
		code = e.ERROR
		global.Log.Warn("OnOrClose failed", err.Error())
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

// GetByIdWithDish @GetByIdWithDish 根据套餐id获取套餐和关联菜品信息
// @Tags setmeal
// @Security JWTAuth
// @Produce json
// @Param id path string true
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/setmeal/{id} [get]
func (sc *SetMealController) GetByIdWithDish(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		id   string
		res  response.SetMealWithDishByIdVo
		err  error
	)
	id = ctx.Param("id")

	uint64Id, _ := strconv.ParseUint(id, 10, 64)
	if res, err = sc.service.GetByIdWithDish(ctx, uint64Id); err != nil {
		code = e.ERROR
		global.Log.Warn("GetByIdWithDish failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Data: nil,
			Msg:  e.GetMsg(code),
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: res,
		Msg:  e.GetMsg(code),
	})
}
