package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/user/request"
	"takeout/internal/api/user/response"
	"takeout/internal/model"
	"takeout/internal/service"
)

type OrderController struct {
	service service.IOrderService
}

func NewOrderController(service service.IOrderService) *OrderController {
	return &OrderController{service: service}
}

// OrderSubmit 用户下单
func (c OrderController) OrderSubmit(ctx *gin.Context) {
	var (
		data    request.OrderSubmitDTO // 创建接收数据模型
		code    = e.SUCCESS
		err     error
		orderVO response.OrderSubmitVO
	)

	if err = ctx.ShouldBindJSON(&data); err != nil {
		global.Log.Debug("OrderSubmit bind param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}
	// 日志打印
	global.Log.Info("用户下单: ", data)

	// 调用service层进行处理
	if orderVO, err = c.service.OrderSubmit(ctx, data); err != nil {
		code = e.ERROR
		global.Log.Debug("OrderSubmit error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Data: nil,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: orderVO,
		Msg:  e.GetMsg(code),
	})
}

// OrderPayment 订单支付
func (c OrderController) OrderPayment(ctx *gin.Context) {
	var (
		data      request.OrderPaymentDTO // 创建接收数据模型
		paymentVO response.OrderPaymentVO
		code      = e.SUCCESS
		err       error
	)

	if err = ctx.ShouldBindJSON(&data); err != nil {
		global.Log.Debug("OrderSubmit bind param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}
	// 日志打印
	global.Log.Info("订单支付: ", data)

	// 调用service层进行处理
	paymentVO = c.service.OrderPayment(ctx, data)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: paymentVO,
		Msg:  e.GetMsg(code),
	})
}

// OrderDetail 根据订单id查询订单详情
func (c OrderController) OrderDetail(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		err     error
		orderVO response.OrderVO
	)
	orderId := ctx.Param("id")
	// 日志打印
	global.Log.Info("根据订单id查询订单详情: ", orderId)

	// 调用service层进行处理
	if orderVO, err = c.service.OrderDetail(ctx, orderId); err != nil {
		code = e.ERROR
		global.Log.Debug("OrderDetail error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Data: nil,
			Msg:  e.GetMsg(code),
		})
		return
	}

	// 数据为空
	if reflect.DeepEqual(orderVO.Order, model.Order{}) {
		ctx.JSON(http.StatusOK, common.Result{
			Code: code,
			Data: nil,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: orderVO,
		Msg:  e.GetMsg(code),
	})

}

// HistoryOrders 查询历史订单
func (c OrderController) HistoryOrders(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		pageDTO request.PageQueryOrderDTO
		err     error
		pageVo  *common.PageResult
	)
	// 获取分页查询参数
	if err = ctx.ShouldBindQuery(&pageDTO); err != nil {
		global.Log.Debug("HistoryOrders bind param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}
	if pageVo, err = c.service.HistoryOrders(ctx, pageDTO); err != nil {
		code = e.ERROR
		global.Log.Debug("HistoryOrders error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Data: nil,
			Msg:  e.GetMsg(code),
		})
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: pageVo,
		Msg:  e.GetMsg(code),
	})
}

// CancelOrder 用户取消订单
func (c OrderController) CancelOrder(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
	)
	orderId := ctx.Param("id")

	// 日志打印
	global.Log.Info("取消订单: ", orderId)

	// 调用service层进行处理
	if err = c.service.CancelOrder(ctx, orderId); err != nil {
		code = e.ERROR
		global.Log.Debug("CancelOrder error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})

}

// RepetitionOrder 再来一单
func (c OrderController) RepetitionOrder(ctx *gin.Context) {

	var (
		code = e.SUCCESS
		err  error
	)
	orderId := ctx.Param("id")

	// 日志打印
	log.Printf("再来一单: [%v]", orderId)

	// 调用service层进行处理
	if err = c.service.RepetitionOrder(ctx, orderId); err != nil {
		code = e.ERROR
		global.Log.Debug("RepetitionOrder error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})

}

// OrderReminder 用户催单
func (c OrderController) OrderReminder(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
	)
	// 接收路径参数
	orderId := ctx.Param("id")
	// 日志打印
	log.Printf("用户催单: [%v]", orderId)
	// 调用service层进行处理
	err = c.service.OrderReminder(ctx, orderId)
	if err != nil {
		code = e.ERROR
		global.Log.Debug("OrderReminder error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}
