package user

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/user/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type DishRouter struct{}

func (dr *DishRouter) InitApiRouter(parent *gin.RouterGroup) {
	privateRouter := parent.Group("dish")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTUser())
	// 依赖注入
	dishCtrl := controller.NewDishController(
		service.NewDishService(dao.NewDishRepo(global.DB), dao.NewDishFlavorDao()),
	)
	{
		privateRouter.GET("/list", dishCtrl.List)
	}
}
