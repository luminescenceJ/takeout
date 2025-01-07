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

type DishController struct {
	service service.IDishService
}

// AddDish @AddDish 新增菜品数据接口
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param data body request.DishDTO true "新增信息"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/dish [post]
func (c DishController) AddDish(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		err     error
		dishDTO request.DishDTO
	)
	if err = ctx.Bind(&dishDTO); err != nil {
		global.Log.Debug("param DishDTO failed")
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = c.service.AddDishWithFlavors(ctx, dishDTO); err != nil {
		code = e.ERROR
		global.Log.Warn("AddDish failed", err.Error())
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// PageQuery @PageQuery 菜品分页查询
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param data body request.DishPageQueryDTO true "新增信息"
// @Success 200 {object} common.Result{data=common.PageResult} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/dish/page [get]
func (c DishController) PageQuery(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		err     error
		pageDTO request.DishPageQueryDTO
		pageRES *common.PageResult
	)

	if err = ctx.Bind(&pageDTO); err != nil {
		global.Log.Debug("param DishPageQueryDTO failed")
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if pageRES, err = c.service.PageQuery(ctx, &pageDTO); err != nil {
		code = e.ERROR
		global.Log.Warn("PageQuery failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: pageRES,
		Msg:  e.GetMsg(code),
	})
}

// GetById @GetById 根据id查询菜品信息
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param id path string true "菜品id"
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/dish/{id} [get]
func (c DishController) GetById(ctx *gin.Context) {
	var (
		code   = e.SUCCESS
		err    error
		id     uint64
		dishVO response.DishVo
	)
	id, _ = strconv.ParseUint(ctx.Param("id"), 10, 64)
	if dishVO, err = c.service.GetByIdWithFlavors(ctx, id); err != nil {
		code = e.ERROR
		global.Log.Warn("GetById failed : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: dishVO,
		Msg:  e.GetMsg(code),
	})
}

// List @List 根据分类id查询菜品信息
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param categoryId query string true "分类id"
// @Success 200 {object} common.Result{data:response.DishListVo} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/dish/list [get]
func (c DishController) List(ctx *gin.Context) {
	var (
		code       = e.SUCCESS
		err        error
		categoryId uint64
		dishList   []response.DishListVo
	)

	categoryId, err = strconv.ParseUint(ctx.Query("categoryId"), 10, 64)

	if dishList, err = c.service.List(ctx, categoryId); err != nil {
		code = e.ERROR
		global.Log.Warn("List failed : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: dishList,
		Msg:  e.GetMsg(code),
	})
}

// Update @Update 修改菜品信息
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param dto body request.DishUpdateDTO true
// @Success 200 {object} common.Result{} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/dish [put]
func (c DishController) Update(ctx *gin.Context) {
	var (
		code      = e.SUCCESS
		err       error
		updateDTO request.DishUpdateDTO
	)
	if err = ctx.ShouldBindJSON(&updateDTO); err != nil {
		global.Log.Debug("param Update failed")
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = c.service.Update(ctx, updateDTO); err != nil {
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

// OnOrClose @OnOrClose 菜品启售或禁售
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param id query string true
// @Param status path string true
// @Success 200 {object} common.Result{} "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/dish/status/{status} [post]
func (c DishController) OnOrClose(ctx *gin.Context) {
	var (
		code   = e.SUCCESS
		err    error
		id     uint64
		status int
	)
	id, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
	status, _ = strconv.Atoi(ctx.Param("status"))

	if err = c.service.OnOrClose(ctx, id, status); err != nil {
		code = e.ERROR
		global.Log.Warn("OnOrClose failed: ", err.Error())
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

// Delete @Delete 删除菜品信息
// @Tags dish
// @Security JWTAuth
// @Produce json
// @Param ids query string true
// @Success 200 {object} common.Result{} "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/dish [delete]
func (c DishController) Delete(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
		ids  string
	)
	ids = ctx.Query("ids")

	if err = c.service.Delete(ctx, ids); err != nil {
		code = e.ERROR
		global.Log.Warn("Delete failed: ", err.Error())
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

func NewDishController(service service.IDishService) *DishController {
	return &DishController{service: service}
}
