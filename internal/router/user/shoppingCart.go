package user

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/user/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type ShoppingCartRouter struct{}

func (dr *ShoppingCartRouter) InitApiRouter(parent *gin.RouterGroup) {
	privateRouter := parent.Group("shoppingCart")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTUser())
	// 依赖注入

	shoppingCartCtrl := controller.NewShoppingController(
		service.NewShoppingCartService(dao.NewShoppingCartDao(global.DB)),
	)

	{
		//privateRouter.GET("/list", dishCtrl.List)

		// 添加购物车
		privateRouter.POST("add", shoppingCartCtrl.AddShoppingCart)
		// 查看购物车
		privateRouter.GET("list", shoppingCartCtrl.QueryShoppingCart)
		// 清空购物车
		privateRouter.DELETE("clean", shoppingCartCtrl.CleanShoppingCart)
		// 减少购物车某项
		privateRouter.POST("sub", shoppingCartCtrl.SubShoppingCart)
	}
}
