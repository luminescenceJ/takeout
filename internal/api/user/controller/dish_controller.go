package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/admin/response"
	"takeout/internal/service"
)

type DishController struct {
	service service.IDishService
}

func NewDishController(service service.IDishService) *DishController {
	return &DishController{service: service}
}

// List @List 根据分类id查询菜品信息
// @Tags Userdish
// @Security JWTAuth
// @Produce json
// @Param categoryId query string true "分类id"
// @Success 200 {object} common.Result{Data=response.DishListVo} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/dish/list [get]
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
