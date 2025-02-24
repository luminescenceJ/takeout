package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	iUtils "github.com/iWyh2/go-myUtils/utils"
	"github.com/robfig/cron/v3"
	"github.com/ulule/deepcopier"
	"log"
	"strconv"
	"takeout/common"
	"takeout/common/enum"
	"takeout/global"
	"takeout/internal/api/user/request"
	"takeout/internal/api/user/response"
	"takeout/internal/model"
	"takeout/internal/router/websocket"
	"takeout/repository"
	"time"
)

type IOrderService interface {
	OrderReminder(ctx *gin.Context, orderId string) error
	RepetitionOrder(ctx *gin.Context, orderId string) error
	OrderSubmit(ctx *gin.Context, data request.OrderSubmitDTO) (response.OrderSubmitVO, error)
	OrderPayment(ctx *gin.Context, orderData request.OrderPaymentDTO) response.OrderPaymentVO
	CancelOrder(ctx *gin.Context, orderId string) error
	HistoryOrders(ctx *gin.Context, dto request.PageQueryOrderDTO) (*common.PageResult, error)

	OrderDetail(ctx *gin.Context, orderId string) (response.OrderVO, error)

	OrderConfirm(ctx *gin.Context, data request.OrderConfirmDTO) error
	OrderRejection(ctx *gin.Context, data request.OrderRejectionDTO) error
	OrderDelivery(ctx *gin.Context, orderId string) error
	OrderComplete(ctx *gin.Context, orderId string) error
	CancelOrderByBusiness(ctx *gin.Context, data request.OrderCancelDTO) error
	OrderConditionSearch(ctx *gin.Context, data request.OrderPageQueryDTO) (*common.PageResult, error)
	OrderStatistics(ctx *gin.Context) (response.OrderStatisticsVO, error)
}
type OrderService struct {
	repo repository.OrderRepo
}

func NewOrderService(repo repository.OrderRepo) IOrderService {
	service := &OrderService{repo: repo}
	defer func() {
		global.Log.Info("启动定时器: [%s]", time.Now().Format("2006-01-02 15:04:05"))
		// 获得定时器
		timerTask := cron.New(cron.WithSeconds())
		// 添加定时器任务
		if _, err := timerTask.AddFunc("0 * * * * ?", service.processTimeoutOrder); err != nil {
			global.Log.Warn("TimerTaskError")
		}
		//if _, err := timerTask.AddFunc("0 0 1 * * ?", service.processDeliveryOrder); err != nil {
		//	panic(errs.TimerTaskError)
		//}
		// 启动定时器
		timerTask.Start()
	}()
	return service
}

// OrderSubmit 用户下单
func (s *OrderService) OrderSubmit(ctx *gin.Context, data request.OrderSubmitDTO) (response.OrderSubmitVO, error) {
	userId, ok := ctx.Get(enum.CurrentId)
	if !ok {
		return response.OrderSubmitVO{}, errors.New("未查找到用户")
	}
	OrderVo, err := s.repo.OrderSubmit(ctx, data, int(userId.(uint64)))
	if err != nil {
		return response.OrderSubmitVO{}, err
	}
	return OrderVo, nil
}

// OrderPayment 订单支付
func (s *OrderService) OrderPayment(ctx *gin.Context, orderData request.OrderPaymentDTO) response.OrderPaymentVO {
	// 调用微信支付接口，生成预支付交易单，此处进行模拟
	log.Printf("调用微信支付接口: [%v]", orderData)
	// 模拟支付成功，修改订单状态
	go s.paySuccess(ctx, orderData.OrderNumber)
	return response.OrderPaymentVO{
		NonceStr:   iUtils.UUID(),
		PaySign:    "111",
		SignType:   "111",
		PackageStr: iUtils.UUID(),
		TimeStamp:  strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}

// OrderDetail 根据订单id查询订单详情
func (s *OrderService) OrderDetail(ctx *gin.Context, orderId string) (response.OrderVO, error) {
	var (
		order           *model.Order
		orderDetailList []model.OrderDetail
		err             error
		orderVO         response.OrderVO
	)
	// 根据id查询订单
	order, err = s.repo.GetOrderById(orderId)
	if err != nil {
		return response.OrderVO{}, err
	}
	// 查询该订单对应的菜品/套餐明细
	orderDetailList, err = s.repo.GetOrderDetailByOrderId(orderId)
	if err != nil {
		return response.OrderVO{}, err
	}
	// 将该订单及其详情封装到OrderVO并返回
	err = deepcopier.Copy(order).To(&orderVO)
	if err != nil {
		return response.OrderVO{}, err
	}
	orderVO.OrderDetailList = orderDetailList

	return orderVO, nil
}

// CancelOrder 用户取消订单
func (s *OrderService) CancelOrder(ctx *gin.Context, orderId string) error {
	// 根据id查询订单
	order, _ := s.repo.GetOrderById(orderId)
	// 校验订单是否存在 ,校验订单状态
	if order == nil || order.Status > enum.ToBeConfirmed {
		return errors.New("订单错误")
	}

	// 订单处于待接单状态下取消，需要进行退款
	if order.Status == enum.ToBeConfirmed {
		// 模拟微信退款
		log.Printf("待接单订单取消, 退款: [%v￥]", order.Amount)
		// 支付状态修改为 退款
		order.PayStatus = enum.Refund
	}

	// 更新订单状态、取消原因、取消时间
	order.Status = enum.Cancelled
	order.CancelReason = "用户取消"
	order.CancelTime = model.LocalTime(time.Now())
	err := s.repo.UpdateOrder(order)
	if err != nil {
		return err
	}
	return nil
}

// RepetitionOrder 再来一单
func (s *OrderService) RepetitionOrder(ctx *gin.Context, orderId string) error {
	// 查询当前用户id
	userId := int(ctx.MustGet(enum.CurrentId).(uint64))
	// 根据订单id查询当前订单详情
	return s.repo.RepetitionOrder(ctx, orderId, userId)
}

// OrderReminder 用户催单
func (s *OrderService) OrderReminder(ctx *gin.Context, orderId string) error {
	// 查询订单是否存在
	order, _ := s.repo.GetOrderById(orderId)
	if order == nil {
		return errors.New("订单不存在")
	}
	// 基于WebSocket实现催单
	jsonMap := map[string]any{
		"type":    2,
		"orderId": orderId,
		"content": "订单号: " + order.Number,
	}
	websocket.WSServer.SendToAllClients(jsonMap)
	return nil

}

// HistoryOrders 查询历史订单
func (s *OrderService) HistoryOrders(ctx *gin.Context, dto request.PageQueryOrderDTO) (*common.PageResult, error) {
	// 查询当前用户id
	userId := int(ctx.MustGet(enum.CurrentId).(uint64))
	return s.repo.HistoryOrders(ctx, dto, userId)
}

// 支付成功，修改订单状态
func (s *OrderService) paySuccess(ctx *gin.Context, orderNumber string) {
	userId, _ := ctx.Get(enum.CurrentId)
	// 根据订单号查询当前用户的订单
	order, _ := s.repo.GetOrderByNumberAndUserId(orderNumber, strconv.FormatUint(userId.(uint64), 10))
	// 根据订单id更新订单的状态、支付方式、支付状态、结账时间
	err := s.repo.UpdateOrder(&model.Order{
		Id:           order.Id,
		Status:       enum.ToBeConfirmed,
		PayStatus:    enum.Paid,
		CheckoutTime: model.LocalTime(time.Now()),
	})
	if err != nil {
		return
	}
	// 基于WebSocket提醒商家来单了
	jsonMap := map[string]any{
		"type":    1,
		"orderId": order.Id,
		"content": "订单号: " + orderNumber,
	}
	websocket.WSServer.SendToAllClients(jsonMap)
}

// 处理支付超时订单
func (s *OrderService) processTimeoutOrder() {
	// 当前时间
	nowTime := time.Now()
	log.Printf("处理支付超时订单: [%s]", time.Now().Format("2006-01-02 15:04:05"))
	// 15分钟前的时间
	duration, _ := time.ParseDuration("-1m")
	nowTime.Add(15 * duration)
	// 查询出超时订单
	ordersList, _ := s.repo.GetOrderByStatusAndOrderTime(enum.PendingPayment, model.LocalTime(nowTime))

	if ordersList != nil && len(ordersList) > 0 {
		for _, order := range ordersList {
			order.Status = enum.Cancelled
			order.CancelReason = "支付超时，自动取消"
			order.CancelTime = model.LocalTime(time.Now())
			_ = s.repo.UpdateOrder(&order)
		}
	}
}

// OrderConfirm 接单
func (s *OrderService) OrderConfirm(ctx *gin.Context, data request.OrderConfirmDTO) error {
	err := s.repo.OrderConfirm(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// OrderRejection 拒绝订单
func (s *OrderService) OrderRejection(ctx *gin.Context, data request.OrderRejectionDTO) error {
	err := s.repo.OrderRejection(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// CancelOrderByBusiness 商家取消订单
func (s *OrderService) CancelOrderByBusiness(ctx *gin.Context, data request.OrderCancelDTO) error {
	order, err := s.repo.GetOrderById(strconv.Itoa(data.OrderId))
	if err != nil {
		return err
	}
	if order.PayStatus == enum.Paid {
		// 用户已支付，需要退款
		// 模拟微信退款
		log.Printf("已支付订单被取消订单, 退款: [%v￥]", order.Amount)
	}
	// 根据订单id更新订单状态、取消原因、取消时间
	order.Status = enum.Cancelled
	order.CancelReason = data.CancelReason
	order.CancelTime = model.LocalTime(time.Now())
	err = s.repo.UpdateOrder(order)
	if err != nil {
		return err
	}
	return nil
}

// OrderDelivery 订单派送
func (s *OrderService) OrderDelivery(ctx *gin.Context, orderId string) error {
	err := s.repo.OrderDelivery(orderId)
	if err != nil {
		return err
	}
	return nil
}

// OrderComplete 完成订单
func (s *OrderService) OrderComplete(ctx *gin.Context, orderId string) error {
	err := s.repo.OrderComplete(orderId)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) OrderConditionSearch(ctx *gin.Context, data request.OrderPageQueryDTO) (*common.PageResult, error) {
	return s.repo.OrderConditionSearch(ctx, data)
}

func (s *OrderService) OrderStatistics(ctx *gin.Context) (response.OrderStatisticsVO, error) {
	return s.repo.OrderStatistics(ctx)
}
