package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/api/controller"
	"takeout/middle"
)

type CommonRouter struct{}

func (dr *CommonRouter) InitApiRouter(parent *gin.RouterGroup) {
	privateRouter := parent.Group("common")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTAdmin())
	commCtrl := new(controller.CommonController)
	{
		privateRouter.POST("/upload", commCtrl.Upload)
	}
}
