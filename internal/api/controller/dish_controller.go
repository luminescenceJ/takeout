package controller

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/service"
)

type DishController struct {
	service service.IDishService
}

// AddDish 新增菜品数据
func (c DishController) AddDish(context *gin.Context) {

}

// PageQuery 菜品分页查询
func (c DishController) PageQuery(context *gin.Context) {

}

// GetById 根据id查询菜品信息
func (c DishController) GetById(context *gin.Context) {

}

// List 根据分类id查询菜品信息
func (c DishController) List(context *gin.Context) {

}

// Update 修改菜品信息
func (c DishController) Update(context *gin.Context) {

}

// OnOrClose 菜品启售或禁售
func (c DishController) OnOrClose(context *gin.Context) {

}

// Delete 删除菜品信息
func (c DishController) Delete(context *gin.Context) {

}

func NewDishController(service service.IDishService) *DishController {
	return &DishController{service: service}
}
