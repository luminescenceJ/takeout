package repository

import (
	"context"
	"takeout/common"
	"takeout/internal/api/user/request"
	"takeout/internal/api/user/response"
	"takeout/internal/model"
)

type OrderRepo interface {
	GetOrderDetailByOrderId(orderId string) (orderDetailList []model.OrderDetail, err error)
	OrderSubmit(ctx context.Context, data request.OrderSubmitDTO, userId int) (response.OrderSubmitVO, error)
	RepetitionOrder(ctx context.Context, orderId string, userId int) error
	OrderReminder(orderId string)
	GetOrderById(orderId string) (*model.Order, error)
	UpdateOrder(order *model.Order) error
	GetOrderByNumberAndUserId(orderNumber, userId string) (*model.Order, error)
	HistoryOrders(ctx context.Context, data request.PageQueryOrderDTO, userId int) (*common.PageResult, error)

	//OrderConfirm(ctx context.Context, data request.OrderConfirmDTO) error
	//OrderRejection(ctx context.Context, data *request.OrderRejectionDTO) error
	//OrderDelivery(orderId string) error
	//OrderComplete(orderId string) error
}
