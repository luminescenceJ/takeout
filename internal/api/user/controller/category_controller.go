package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/service"
)

type CategoryController struct {
	service service.ICategoryService
}

func NewCategoryController(service service.ICategoryService) *CategoryController {
	return &CategoryController{service: service}
}

// List @List 根据类型查询分类
// @Tags UserCategory
// @Security JWTAuth
// @Produce json
// @Param type query string true "查询信息"
// @Success 200 {object} common.Result "success"
// @Router /user/category/list [get]
func (cc *CategoryController) List(ctx *gin.Context) {
	var code = e.SUCCESS
	cate, _ := strconv.Atoi(ctx.Query("type"))
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
