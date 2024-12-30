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

type CategoryController struct {
	service service.ICategoryService
}

func NewCategoryController(service service.ICategoryService) *CategoryController {
	return &CategoryController{service: service}
}

// @AddCategory 新增分类接口
// @Tags Category
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
	err = ctx.ShouldBindBodyWithJSON(&categoryDTO)
	if err != nil {
		global.Log.Debug("param CategoryDTO json failed", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
	}
	if err = cc.service.AddCategory(ctx, categoryDTO); err != nil {
		code = e.ERROR
		global.Log.Debug("AddCategory failed", err.Error())
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
	})
}

func (cc *CategoryController) PageQuery(ctx *gin.Context) {
	code := e.SUCCESS

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
	})
}

func (cc *CategoryController) List(ctx *gin.Context) {
	code := e.SUCCESS

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
	})
}

func (cc *CategoryController) DeleteById(ctx *gin.Context) {
	code := e.SUCCESS

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
	})
}

func (cc *CategoryController) EditCategory(ctx *gin.Context) {
	code := e.SUCCESS

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
	})
}

func (cc *CategoryController) SetStatus(ctx *gin.Context) {
	code := e.SUCCESS

	ctx.JSON(http.StatusOK, common.Result{Code: code, Msg: e.GetMsg(code)})
}
