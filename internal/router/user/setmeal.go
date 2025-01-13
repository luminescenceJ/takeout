package user

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/user/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type SetmealRouter struct{}

func (dr *SetmealRouter) InitApiRouter(parent *gin.RouterGroup) {
	privateRouter := parent.Group("setmeal")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTUser())
	// 依赖注入
	dishCtrl := controller.NewSetMealController(
		service.NewSetMealService(dao.NewSetMealDao(global.DB), dao.NewSetMealDishDao()),
	)
	{
		privateRouter.GET("/dish/:id", dishCtrl.GetDishByCategoryId)
		privateRouter.GET("/list", dishCtrl.List)
	}
}
