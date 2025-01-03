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

type DishController struct {
	service service.IDishService
}

// @AddDish 新增菜品数据接口
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

// PageQuery 菜品分页查询
func (c DishController) PageQuery(ctx *gin.Context) {

}

// GetById 根据id查询菜品信息
func (c DishController) GetById(ctx *gin.Context) {

}

// List 根据分类id查询菜品信息
func (c DishController) List(ctx *gin.Context) {

}

// Update 修改菜品信息
func (c DishController) Update(ctx *gin.Context) {

}

// OnOrClose 菜品启售或禁售
func (c DishController) OnOrClose(ctx *gin.Context) {

}

// Delete 删除菜品信息
func (c DishController) Delete(ctx *gin.Context) {

}

func NewDishController(service service.IDishService) *DishController {
	return &DishController{service: service}
}
