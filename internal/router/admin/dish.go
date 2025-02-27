package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/admin/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type DishRouter struct{}

func (dr *DishRouter) InitApiRouter(parent *gin.RouterGroup) {
	//publicRouter := parent.Group("category")
	privateRouter := parent.Group("dish")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTAdmin())
	// 依赖注入
	dishCtrl := controller.NewDishController(
		service.NewDishService(dao.NewDishRepo(global.DB), dao.NewDishFlavorDao()),
	)
	{
		privateRouter.POST("", dishCtrl.AddDish)
		privateRouter.GET("/page", dishCtrl.PageQuery)
		privateRouter.GET("/:id", dishCtrl.GetById)
		privateRouter.GET("/list", dishCtrl.List)
		privateRouter.PUT("", dishCtrl.Update)

		privateRouter.POST("/status/:status", dishCtrl.OnOrClose)
		privateRouter.DELETE("", dishCtrl.Delete)
	}
}
