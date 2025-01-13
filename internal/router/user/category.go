package user

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/user/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type CategoryRouter struct {
	service service.ICategoryService
}

func (cr *CategoryRouter) InitApiRouter(router *gin.RouterGroup) {

	privateRouter := router.Group("category")

	privateRouter.Use(middle.VerifiyJWTUser()) // 私有路由使用jwt验证

	//依赖注入
	cr.service = service.NewCategoryService(dao.NewCategoryDao(global.DB))
	wxCtrl := controller.NewCategoryController(cr.service)

	{
		privateRouter.GET("list", wxCtrl.List)
	}
}
