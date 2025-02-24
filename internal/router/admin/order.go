package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/api/admin/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type OrderRouter struct {
	service service.IOrderService
}

func (er *OrderRouter) InitApiRouter(router *gin.RouterGroup) {
	// /admin/order
	//publicRouter := router.Group("order")
	privateRouter := router.Group("order")

	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTAdmin())

	// 依赖注入
	er.service = service.NewOrderService(dao.NewOrderDao())
	orderCtl := controller.NewOrderController(er.service)
	{
		// 接单
		privateRouter.PUT("confirm", orderCtl.OrderConfirm)
		// 拒单
		privateRouter.PUT("rejection", orderCtl.OrderRejection)
		// 商家取消订单
		privateRouter.PUT("cancel", orderCtl.CancelOrder)
		// 派送订单
		privateRouter.PUT("delivery/:id", orderCtl.OrderDelivery)
		// 完成订单
		privateRouter.PUT("complete/:id", orderCtl.OrderComplete)
		// 订单搜索
		privateRouter.GET("conditionSearch", orderCtl.OrderConditionSearch)
		// 查询订单详情
		privateRouter.GET("details/:id", orderCtl.OrderDetail)
		// 各个状态的订单数量统计
		privateRouter.GET("statistics", orderCtl.OrderStatistics)
	}
}
