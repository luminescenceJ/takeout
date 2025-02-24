package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/user/request"
	"takeout/internal/api/user/response"
	"takeout/internal/service"
)

type OrderController struct {
	service service.IOrderService
}

func NewOrderController(orderService service.IOrderService) *OrderController {
	return &OrderController{service: orderService}
}

// OrderConfirm 接单
func (c *OrderController) OrderConfirm(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
		data request.OrderConfirmDTO
	)
	// 确保请求体是JSON格式
	if err = ctx.ShouldBindJSON(&data); err != nil {
		code = e.ERROR
		global.Log.Debug("Invalid request payload:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  "Invalid request payload",
		})
		return
	}
	// 日志打印
	log.Printf("接单: [%v]", data.OrderId)
	// 调用service层进行处理

	if err = c.service.OrderConfirm(ctx, data); err != nil {
		code = e.ERROR
		global.Log.Warn("OrderConfirm Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  "wrong password or error on json web token generate",
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// OrderRejection 拒单
func (c *OrderController) OrderRejection(ctx *gin.Context) {

	var (
		code = e.SUCCESS
		err  error
		data request.OrderRejectionDTO
	)
	// 确保请求体是JSON格式
	if err = ctx.ShouldBindJSON(&data); err != nil {
		code = e.ERROR
		global.Log.Debug("Invalid request payload:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  "Invalid request payload",
		})
		return
	}
	// 日志打印
	log.Printf("拒单: [%v]", data.OrderId)
	// 调用service层进行处理

	if err = c.service.OrderRejection(ctx, data); err != nil {
		code = e.ERROR
		global.Log.Warn("OrderRejection Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  "wrong password or error on json web token generate",
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})

}

// CancelOrder 商家取消订单
func (c *OrderController) CancelOrder(ctx *gin.Context) {

	var (
		code = e.SUCCESS
		err  error
		data request.OrderCancelDTO
	)
	// 确保请求体是JSON格式
	if err = ctx.ShouldBindJSON(&data); err != nil {
		code = e.ERROR
		global.Log.Debug("Invalid request payload:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  "Invalid request payload",
		})
		return
	}
	// 日志打印
	log.Printf("商家取消订单: [%v]", data.OrderId)
	// 调用service层进行处理
	if err = c.service.CancelOrderByBusiness(ctx, data); err != nil {
		code = e.ERROR
		global.Log.Warn("CancelOrderByBusiness Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  "wrong password or error on json web token generate",
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// OrderDelivery 订单派送
func (c *OrderController) OrderDelivery(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
	)

	// 接收路径参数
	orderId := ctx.Param("id")
	// 日志打印
	log.Printf("订单派送: [%v]", orderId)
	// 调用service层进行处理
	if err = c.service.OrderDelivery(ctx, orderId); err != nil {
		code = e.ERROR
		global.Log.Warn("OrderDelivery Error:", err.Error())
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

// OrderComplete 完成订单
func (c *OrderController) OrderComplete(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
	)

	// 接收路径参数
	orderId := ctx.Param("id")
	// 日志打印
	log.Printf("完成订单: [%v]", orderId)
	// 调用service层进行处理
	if err = c.service.OrderComplete(ctx, orderId); err != nil {
		code = e.ERROR
		global.Log.Warn("OrderDelivery Error:", err.Error())
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

// OrderConditionSearch 订单搜索
func (c *OrderController) OrderConditionSearch(ctx *gin.Context) {
	var (
		code       = e.SUCCESS
		err        error
		data       request.OrderPageQueryDTO
		pageResult *common.PageResult
	)
	// 确保请求体是JSON格式
	if err = ctx.ShouldBindQuery(&data); err != nil {
		code = e.ERROR
		global.Log.Debug("Invalid request payload:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{
			Code: code,
			Msg:  "Invalid request payload",
		})
		return
	}
	// 日志打印
	log.Printf("订单搜索: %v", data)
	// 调用service层进行处理
	if pageResult, err = c.service.OrderConditionSearch(ctx, data); err != nil {
		code = e.ERROR
		global.Log.Warn("OrderConditionSearch Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: pageResult,
		Msg:  e.GetMsg(code),
	})

}

// OrderStatistics 各个状态的订单数量统计
func (c *OrderController) OrderStatistics(ctx *gin.Context) {

	var (
		code   = e.SUCCESS
		err    error
		result response.OrderStatisticsVO
	)

	// 调用service层进行处理
	if result, err = c.service.OrderStatistics(ctx); err != nil {
		code = e.ERROR
		global.Log.Warn("OrderConditionSearch Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	log.Printf("各个状态的订单数量统计: %v", result)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: result,
		Msg:  e.GetMsg(code),
	})
}

// OrderDetail 查询订单详情
func (c *OrderController) OrderDetail(ctx *gin.Context) {
	var (
		result response.OrderVO
		code   = e.SUCCESS
		err    error
	)

	// 接收路径参数
	orderId := ctx.Param("id")
	// 日志打印
	log.Printf("查询订单详情: [%v]", orderId)
	// 调用service层进行处理
	if result, err = c.service.OrderDetail(ctx, orderId); err != nil {
		code = e.ERROR
		global.Log.Warn("OrderDelivery Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: result,
		Msg:  e.GetMsg(code),
	})
}
