package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type CategoryRouter struct {
	service service.ICategoryService
}

func (cr *CategoryRouter) InitApiRouter(router *gin.RouterGroup) {
	// admin/category
	privateRouter := router.Group("category")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTAdmin())
	//依赖注入
	cr.service = service.NewCategoryService(dao.NewCategoryDao(global.DB))
	categoryCtrl := controller.NewCategoryController(cr.service)

	{
		privateRouter.POST("", categoryCtrl.AddCategory)
		privateRouter.GET("page", categoryCtrl.PageQuery)
		privateRouter.GET("list", categoryCtrl.List)
		privateRouter.DELETE("", categoryCtrl.DeleteById)
		privateRouter.PUT("", categoryCtrl.EditCategory)
		privateRouter.POST("status/:status", categoryCtrl.SetStatus)
	}
}
