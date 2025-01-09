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
// @Tags setmeal
// @Security JWTAuth
// @Produce json
// @Param dto body request.SetMealDTO true "菜品信息"
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
		global.Log.Debug("SaveWithDish保存套餐和菜品信息 结构体解析失败:", err.Error())
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
// @Param dto body request.SetMealPageQueryDTO true "分页查询dto"
// @Success 200 {object} common.Result{Data=response.SetMealPageQueryVo} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/setmeal/page [get]
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
// @Param id query string true "id"
// @Param status path string true "status"
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
// @Param id path string true "id"
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

// Update @Update 更新套餐和其关联菜品
// @Tags setmeal
// @Security JWTAuth
// @Produce json
// @Param dto body request.SetMealDTO true "更新套餐和其关联菜品结构体"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/setmeal [put]
func (sc *SetMealController) Update(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
		dto  request.SetMealDTO
	)
	// 这里price字段存在bug，在修改注释的情况下发送的是string类型，其他时候是number类型
	// 需要在json Unmarshal的时候进行判定

	//// 读取并打印请求体的原始 JSON 数据
	//body, err := io.ReadAll(ctx.Request.Body)
	//if err != nil {
	//	global.Log.Warn("读取请求体失败:", err.Error())
	//	return
	//}
	//// 打印请求体的 JSON 数据
	//global.Log.Warn("接收到的请求体 JSON:", string(body))
	//ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err = ctx.ShouldBindJSON(&dto); err != nil {
		global.Log.Warn("Update Invalid request payload:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	if err = sc.service.Update(ctx, dto); err != nil {
		code = e.ERROR
		global.Log.Warn("Update Internal Server Faliure:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Data: nil,
			Msg:  e.GetMsg(code),
		})
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// DeleteBatch @DeleteBatch 批量删除套餐和其关联菜品
// @Tags setmeal
// @Security JWTAuth
// @Produce json
// @Param ids query string true "id集合"
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/setmeal [delete]
func (sc *SetMealController) DeleteBatch(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
		ids  string
	)
	ids = ctx.Query("ids")
	if ids == "" {
		code = e.SUCCESS
		ctx.JSON(http.StatusOK, common.Result{})
		return
	}

	if err = sc.service.DeleteBatch(ctx, ids); err != nil {
		code = e.ERROR
		global.Log.Warn("DeleteBatch Internal Server Faliure:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Data: nil,
			Msg:  e.GetMsg(code),
		})
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}
