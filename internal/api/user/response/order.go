package response

import (
	"takeout/internal/model"
)

// OrderOverViewVO 订单概览数据
type OrderOverViewVO struct {
	WaitingOrders   int `json:"waitingOrders"`
	DeliveredOrders int `json:"deliveredOrders"`
	CompletedOrders int `json:"completedOrders"`
	CancelledOrders int `json:"cancelledOrders"`
	AllOrders       int `json:"allOrders"`
}

// OrderPaymentVO 订单支付返回数据模型
type OrderPaymentVO struct {
	// 随机字符串
	NonceStr string `json:"nonceStr"`
	// 签名
	PaySign string `json:"paySign"`
	// 时间戳
	TimeStamp string `json:"timeStamp"`
	// 签名算法
	SignType string `json:"signType"`
	// 统一下单接口返回的 prepay_id 参数值
	PackageStr string `json:"packageStr"`
}

// OrderReportVO 订单统计返回数据模型
type OrderReportVO struct {
	DateList            string  `json:"dateList"`
	OrderCountList      string  `json:"orderCountList"`
	ValidOrderCountList string  `json:"validOrderCountList"`
	TotalOrderCount     int     `json:"totalOrderCount"`
	ValidOrderCount     int     `json:"validOrderCount"`
	OrderCompletionRate float64 `json:"orderCompletionRate"`
}

// OrderStatisticsVO 订单数量统计返回数据模型
type OrderStatisticsVO struct {
	ToBeConfirmed      int `json:"toBeConfirmed"`
	Confirmed          int `json:"confirmed"`
	DeliveryInProgress int `json:"deliveryInProgress"`
}

// OrderSubmitVO 用户下单接口返回结果
type OrderSubmitVO struct {
	OrderId     int             `json:"id"`
	OrderNumber string          `json:"orderNumber"`
	OrderAmount float64         `json:"orderAmount"`
	OrderTime   model.LocalTime `json:"orderTime"`
}

// OrderVO 查询订单详情返回数据模型
type OrderVO struct {
	model.Order
	OrderDishes     string              `json:"orderDishes"`
	OrderDetailList []model.OrderDetail `json:"orderDetailList"`
}
