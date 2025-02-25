package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/admin/request"
	"takeout/internal/service"
)

type CategoryController struct {
	service service.ICategoryService
}

func NewCategoryController(service service.ICategoryService) *CategoryController {
	return &CategoryController{service: service}
}

// @AddCategory 新增分类接口
// @Tags Category
// @Security JWTAuth
// @Produce json
// @Param data body request.CategoryDTO true "新增分类信息"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Router /admin/category [post]
func (cc *CategoryController) AddCategory(ctx *gin.Context) {
	var (
		code        = e.SUCCESS
		err         error
		categoryDTO request.CategoryDTO
	)
	if err = ctx.ShouldBindJSON(&categoryDTO); err != nil {
		global.Log.Debug("param CategoryDTO json failed", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: "Invalid request payload"})
		return
	}

	if err = cc.service.AddCategory(ctx, categoryDTO); err != nil {
		code = e.ERROR
		global.Log.Debug("AddCategory failed", err.Error())
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// @PageQuery 查询分类分页接口
// @Tags Category
// @Security JWTAuth
// @Produce json
// @Param data query request.CategoryPageQueryDTO true "查询分类信息"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "PageQuery failed"
// @Router /admin/category/page [get]
func (cc *CategoryController) PageQuery(ctx *gin.Context) {
	var (
		code       = e.SUCCESS
		PageQuery  request.CategoryPageQueryDTO
		PageResult *common.PageResult
		err        error
	)
	if err = ctx.ShouldBindQuery(&PageQuery); err != nil {
		global.Log.Debug("param PageQuery json failed", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: "Invalid request payload"})
		return
	}

	PageResult, err = cc.service.PageQuery(ctx, PageQuery)
	if err != nil {
		code = e.ERROR
		global.Log.Debug("PageQuery failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: PageResult,
		Msg:  e.GetMsg(code),
	})
}

// @List 根据类型查询分类
// @Tags Category
// @Security JWTAuth
// @Produce json
// @Param type query string true "查询信息"
// @Success 200 {object} common.Result "success"
// @Router /admin/category/list [get]
func (cc *CategoryController) List(ctx *gin.Context) {
	var code = e.SUCCESS
	cate, err := strconv.Atoi(ctx.DefaultQuery("type", "0"))
	if err != nil {
		code = e.ERROR
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  "Invalid category type",
		})
		return
	}

	res, err := cc.service.List(ctx, cate)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("Category List failed", err.Error())
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: res,
		Msg:  e.GetMsg(code),
	})
}

// @DeleteById 删除分类的接口
// @Tags Category
// @Security JWTAuth
// @Produce json
// @Param id query string true "分类id"
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "DeleteById failed"
// @Router /admin/category [delete]
func (cc *CategoryController) DeleteById(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		id   uint64
		err  error
	)
	id, err = strconv.ParseUint(ctx.DefaultQuery("id", "0"), 10, 64)
	if err != nil || id == 0 {
		code = e.ERROR
		ctx.JSON(http.StatusBadRequest, common.Result{Code: code, Msg: "Invalid category ID"})
		return
	}

	err = cc.service.DeleteById(ctx, id)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("Category DeleteById failed", err.Error())
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

// @EditCategory 编辑分类的接口
// @Tags Category
// @Security JWTAuth
// @Produce json
// @Param data body request.CategoryDTO true "编辑分类的内容"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "deleteById failed"
// @Router /admin/category [put]
func (cc *CategoryController) EditCategory(ctx *gin.Context) {
	var (
		code        = e.SUCCESS
		err         error
		categoryDTO request.CategoryDTO
	)
	if err = ctx.ShouldBindJSON(&categoryDTO); err != nil {
		global.Log.Debug("param CategoryDTO json failed", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Msg: "Invalid request payload"})
		return
	}

	if err = cc.service.Update(ctx, categoryDTO); err != nil {
		code = e.ERROR
		global.Log.Debug("EditCategory failed", err.Error())
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// @SetStatus 启用或禁用分类接口
// @Tags Category
// @Security JWTAuth
// @Produce json
// @Param id query string true "分类id"
// @Param status path string true "状态"
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "SetStatus failed"
// @Router /admin/category/status/{status} [post]
func (cc *CategoryController) SetStatus(ctx *gin.Context) {
	var (
		code   = e.SUCCESS
		status int
		id     uint64
		err    error
	)
	status, err = strconv.Atoi(ctx.Param("status"))
	if err != nil || (status != 0 && status != 1) {
		code = e.ERROR
		ctx.JSON(http.StatusBadRequest, common.Result{Code: code, Msg: "Invalid status value"})
		return
	}

	id, err = strconv.ParseUint(ctx.DefaultQuery("id", "0"), 10, 64)
	if err != nil || id == 0 {
		code = e.ERROR
		ctx.JSON(http.StatusBadRequest, common.Result{Code: code, Msg: "Invalid category ID"})
		return
	}

	err = cc.service.SetStatus(ctx, id, status)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("Category SetStatus failed", err.Error())
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

//package controller
//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strconv"
//	"takeout/common"
//	"takeout/common/e"
//	"takeout/global"
//	"takeout/internal/api/admin/request"
//	"takeout/internal/service"
//)
//
//type CategoryController struct {
//	service service.ICategoryService
//}
//
//func NewCategoryController(service service.ICategoryService) *CategoryController {
//	return &CategoryController{service: service}
//}
//
//// AddCategory @AddCategory 新增分类接口
//// @Tags Category
//// @Security JWTAuth
//// @Produce json
//// @Param data body request.CategoryDTO true "新增分类信息"
//// @Success 200 {object} common.Result "success"
//// @Failure 400 {object} common.Result "Invalid request payload"
//// @Router /admin/category [post]
//func (cc *CategoryController) AddCategory(ctx *gin.Context) {
//	var (
//		code        = e.SUCCESS
//		err         error
//		categoryDTO request.CategoryDTO
//	)
//	err = ctx.ShouldBindBodyWithJSON(&categoryDTO)
//	if err != nil {
//		global.Log.Debug("param CategoryDTO json failed", err.Error())
//		ctx.JSON(http.StatusBadRequest, common.Result{})
//	}
//	if err = cc.service.AddCategory(ctx, categoryDTO); err != nil {
//		code = e.ERROR
//		global.Log.Debug("AddCategory failed", err.Error())
//	}
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//	})
//}
//
//// PageQuery @PageQuery 查询分类分页接口
//// @Tags Category
//// @Security JWTAuth
//// @Produce json
//// @Param data query request.CategoryPageQueryDTO true "查询分类信息"
//// @Success 200 {object} common.Result "success"
//// @Failure 400 {object} common.Result "Invalid request payload"
//// @Failure 500 {object} common.Result "PageQuery failed"
//// @Router /admin/category/page [get]
//func (cc *CategoryController) PageQuery(ctx *gin.Context) {
//	var (
//		code       = e.SUCCESS
//		PageQuery  request.CategoryPageQueryDTO
//		PageResult *common.PageResult
//		err        error
//	)
//	err = ctx.ShouldBindQuery(&PageQuery)
//	if err != nil {
//		global.Log.Debug("param PageQuery json failed", err.Error())
//		ctx.JSON(http.StatusBadRequest, common.Result{})
//	}
//	PageResult, err = cc.service.PageQuery(ctx, PageQuery)
//	if err != nil {
//		code = e.ERROR
//		global.Log.Debug("PageQuery failed", err.Error())
//		ctx.JSON(http.StatusInternalServerError, common.Result{
//			Code: code,
//		})
//	}
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Data: PageResult,
//	})
//}
//
//// @List 根据类型查询分类
//// @Tags Category
//// @Security JWTAuth
//// @Produce json
//// @Param type query string true "查询信息"
//// @Success 200 {object} common.Result "success"
//// @Router /admin/category/list [get]
//func (cc *CategoryController) List(ctx *gin.Context) {
//	var code = e.SUCCESS
//	cate, _ := strconv.Atoi(ctx.Query("type"))
//	res, err := cc.service.List(ctx, cate)
//	if err != nil {
//		code = e.ERROR
//		global.Log.Warn("Category List failed", err.Error())
//	}
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Data: res,
//		Msg:  e.GetMsg(code),
//	})
//}
//
//// @DeleteById 删除分类的接口
//// @Tags Category
//// @Security JWTAuth
//// @Produce json
//// @Param id query string true "分类id"
//// @Success 200 {object} common.Result "success"
//// @Failure 500 {object} common.Result "DeleteById failed"
//// @Router /admin/category [delete]
//func (cc *CategoryController) DeleteById(ctx *gin.Context) {
//	var (
//		code = e.SUCCESS
//		id   uint64
//		err  error
//	)
//	id, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
//	err = cc.service.DeleteById(ctx, id)
//	if err != nil {
//		code = e.ERROR
//		global.Log.Warn("Category DeleteById failed", err.Error())
//		ctx.JSON(http.StatusInternalServerError, common.Result{
//			Code: code,
//			Msg:  e.GetMsg(code),
//		})
//	}
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//	})
//}
//
//// EditCategory @EditCategory 编辑分类的接口
//// @Tags Category
//// @Security JWTAuth
//// @Produce json
//// @Param data body request.CategoryDTO true "编辑分类的内容"
//// @Success 200 {object} common.Result "success"
//// @Failure 400 {object} common.Result "invalid params failed"
//// @Failure 500 {object} common.Result "deleteById failed"
//// @Router /admin/category [put]
//func (cc *CategoryController) EditCategory(ctx *gin.Context) {
//	var (
//		code        = e.SUCCESS
//		err         error
//		categoryDTO request.CategoryDTO
//	)
//	err = ctx.ShouldBindBodyWithJSON(&categoryDTO)
//	if err != nil {
//		global.Log.Debug("param CategoryDTO json failed", err.Error())
//		ctx.JSON(http.StatusBadRequest, common.Result{})
//	}
//	if err = cc.service.Update(ctx, categoryDTO); err != nil {
//		code = e.ERROR
//		global.Log.Debug("AddCategory failed", err.Error())
//	}
//	ctx.JSON(http.StatusOK, common.Result{
//		Code: code,
//		Msg:  e.GetMsg(code),
//	})
//}
//
//// SetStatus @SetStatus 启用或禁用分类接口
//// @Tags Category
//// @Security JWTAuth
//// @Produce json
//// @Param id query string true "分类id"
//// @Param status path string true "状态"
//// @Success 200 {object} common.Result "success"
//// @Failure 500 {object} common.Result "SetStatus failed"
//// @Router /admin/category/status/{status} [post]
//func (cc *CategoryController) SetStatus(ctx *gin.Context) {
//
//	var (
//		code   = e.SUCCESS
//		status int
//		id     uint64
//		err    error
//	)
//	status, _ = strconv.Atoi(ctx.Param("status"))
//	id, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
//	err = cc.service.SetStatus(ctx, id, status)
//	if err != nil {
//		code = e.ERROR
//		global.Log.Warn("Category SetStatus failed", err.Error())
//	}
//	ctx.JSON(http.StatusOK, common.Result{Code: code, Msg: e.GetMsg(code)})
//}
