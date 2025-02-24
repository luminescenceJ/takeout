package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/api/admin/controller"
	"takeout/middle"
	"takeout/repository/dao"
)

type WorkSpaceRouter struct {
}

func (er *WorkSpaceRouter) InitApiRouter(router *gin.RouterGroup) {
	// /admin/workspace
	privateRouter := router.Group("workspace")

	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTAdmin())

	// 依赖注入
	workspaceCtl := controller.NewWorkSpaceController(dao.NewWorkSpaceDao())
	{
		// 今日数据
		privateRouter.GET("businessData", workspaceCtl.TodayBusinessData)
		// 订单管理
		privateRouter.GET("overviewOrders", workspaceCtl.OverviewOrders)
		// 菜品总览
		privateRouter.GET("overviewDishes", workspaceCtl.OverviewDishes)
		// 套餐总览
		privateRouter.GET("overviewSetmeals", workspaceCtl.OverviewSetmeals)

	}
}
