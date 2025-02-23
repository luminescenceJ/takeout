package user

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/api/user/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type OrderRouter struct{}

func (dr *OrderRouter) InitApiRouter(parent *gin.RouterGroup) {
	privateRouter := parent.Group("order")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTUser())
	orderCtrl := controller.NewOrderController(
		service.NewOrderService(dao.NewOrderDao()),
	)
	{
		// 用户下单
		privateRouter.POST("submit", orderCtrl.OrderSubmit)
		// 订单支付
		privateRouter.PUT("payment", orderCtrl.OrderPayment)
		// 根据订单id查询订单详情
		privateRouter.GET("orderDetail/:id", orderCtrl.OrderDetail)

		// 查询历史订单
		privateRouter.GET("historyOrders", orderCtrl.HistoryOrders)
		// 用户取消订单
		privateRouter.PUT("cancel/:id", orderCtrl.CancelOrder)
		// 再来一单
		privateRouter.POST("repetition/:id", orderCtrl.RepetitionOrder)
		// 用户催单
		privateRouter.GET("reminder/:id", orderCtrl.OrderReminder)
	}
}
