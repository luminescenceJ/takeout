package user

import (
	"github.com/gin-gonic/gin"
	"takeout/middle"
)

type ShoppingCartRouter struct{}

func (dr *ShoppingCartRouter) InitApiRouter(parent *gin.RouterGroup) {
	privateRouter := parent.Group("shoppingCart")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTUser())
	// 依赖注入
	//dishCtrl := controller.NewDishController(
	//	service.NewDishService(dao.NewDishRepo(global.DB), dao.NewDishFlavorDao()),
	//)
	{
		//privateRouter.GET("/list", dishCtrl.List)
	}
}
