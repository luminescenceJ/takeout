package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/admin/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type SetMealRouter struct {
	service service.ISetMealService
}

func (er *SetMealRouter) InitApiRouter(router *gin.RouterGroup) {
	// admin/employee
	privateRouter := router.Group("setmeal")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTAdmin())
	// 依赖注入
	er.service = service.NewSetMealService(dao.NewSetMealDao(global.DB), dao.NewSetMealDishDao())
	setmealCtrl := controller.NewSetMealController(er.service)
	{
		privateRouter.POST("", setmealCtrl.SaveWithDish)
		privateRouter.GET("/page", setmealCtrl.PageQuery)
		privateRouter.GET("/:id", setmealCtrl.GetByIdWithDish)
		privateRouter.POST("/status/:status", setmealCtrl.OnOrClose)
		privateRouter.PUT("", setmealCtrl.Update)
		privateRouter.DELETE("", setmealCtrl.DeleteBatch)
	}
}
