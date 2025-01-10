package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/api/admin/controller"
	"takeout/middle"
)

type ShopRouter struct{}

func (s *ShopRouter) InitApiRouter(parent *gin.RouterGroup) {
	privateRouter := parent.Group("shop")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTAdmin())
	shopCtrl := new(controller.ShopController)
	{
		// 设置营业状态
		privateRouter.PUT(":status", shopCtrl.SetShopStatus)
		// 查询营业状态
		privateRouter.GET("status", shopCtrl.GetShopStatus)
	}
}
