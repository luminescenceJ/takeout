package user

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/user/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type WxUserRouter struct {
	service service.IWxUserService
}

func (cr *WxUserRouter) InitApiRouter(router *gin.RouterGroup) {

	privateRouter := router.Group("user")
	publicRouter := router.Group("user")

	privateRouter.Use(middle.VerifiyJWTAdmin()) // 私有路由使用jwt验证

	//依赖注入
	cr.service = service.NewWxUserService(dao.NewWxUserDao(global.DB))
	wxCtrl := controller.NewWxUserController(cr.service)

	{
		publicRouter.POST("login", wxCtrl.Login)
		privateRouter.POST("logout", wxCtrl.Logout)
	}
}
