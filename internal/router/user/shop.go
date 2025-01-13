package user

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/api/user/controller"
	"takeout/internal/service"
)

type ShopRouter struct {
	service service.ShopService
}

func (s *ShopRouter) InitApiRouter(parent *gin.RouterGroup) {
	publicRouter := parent.Group("shop")

	// 私有路由使用jwt验证
	//privateRouter.Use(middle.VerifiyJWTUser())

	shopCtrl := new(controller.ShopController)
	{
		// 查询营业状态
		publicRouter.GET("status", shopCtrl.GetShopStatus)
	}
}
