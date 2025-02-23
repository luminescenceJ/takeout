package dao

import (
	"context"
	"errors"
	"github.com/ulule/deepcopier"
	"gorm.io/gorm"
	"log"
	"strconv"
	"takeout/common"
	"takeout/common/enum"
	"takeout/global"
	"takeout/internal/api/user/request"
	"takeout/internal/api/user/response"
	"takeout/internal/model"
	"takeout/repository"
	"time"
)

type OrderDao struct {
	db *gorm.DB
	AddressBookDao
	ShoppingCartDao
}

func NewOrderDao() repository.OrderRepo {
	return &OrderDao{
		db:              global.DB,
		AddressBookDao:  AddressBookDao{db: global.DB},
		ShoppingCartDao: ShoppingCartDao{db: global.DB},
	}
}

// GetOrderDetailByOrderId 根据订单id查询菜品/套餐明细
func (d OrderDao) GetOrderDetailByOrderId(orderId string) (orderDetailList []model.OrderDetail, err error) {
	if err = global.DB.Where("order_id = ?", orderId).
		Find(&orderDetailList).Error; err != nil {
		return []model.OrderDetail{}, err
	}
	return orderDetailList, nil
}

// OrderSubmit 用户下单
func (d OrderDao) OrderSubmit(ctx context.Context, data request.OrderSubmitDTO, userId int) (response.OrderSubmitVO, error) {
	// 开启事务
	var (
		addressBook      model.AddressBook
		shoppingCartList []model.ShoppingCart
		order            model.Order
		err              error
	)
	if addressBook, err = d.AddressBookDao.GetAddressById(ctx, uint64(data.AddressBookId)); err != nil {
		return response.OrderSubmitVO{}, err
	}
	shoppingCartList = d.ShoppingCartDao.queryShoppingCart(ctx, model.ShoppingCart{UserId: userId})
	if shoppingCartList == nil || len(shoppingCartList) == 0 {
		global.Log.Warn("错误的购物车")
		return response.OrderSubmitVO{}, nil
	}
	_ = deepcopier.Copy(data).To(&order)
	order.UserId = userId
	order.Number = strconv.FormatInt(time.Now().UnixMilli(), 10)
	order.Phone = addressBook.Phone
	order.Address = addressBook.Detail
	order.Consignee = addressBook.Consignee
	order.Status = enum.PendingPayment
	order.PayStatus = enum.UnPaid
	// 向订单表插入1条数据
	if err = global.DB.Create(&order).Error; err != nil {
		return response.OrderSubmitVO{}, err
	}
	// 订单明细数据
	orderDetailList := make([]model.OrderDetail, 0)
	for _, shoppingCart := range shoppingCartList {
		var orderDetail model.OrderDetail
		if err = deepcopier.Copy(&shoppingCart).To(&orderDetail); err != nil {
			return response.OrderSubmitVO{}, err
		}
		orderDetail.OrderId = order.Id
		orderDetailList = append(orderDetailList, orderDetail)
	}
	// 向明细表插入n条数据
	if err = insertBatchOrderDetail(orderDetailList); err != nil {
		return response.OrderSubmitVO{}, err
	}
	// 清理购物车中的数据
	if err = d.Clean(ctx, userId); err != nil {
		return response.OrderSubmitVO{}, err
	}
	// 封装返回结果
	return response.OrderSubmitVO{
		OrderId:     order.Id,
		OrderNumber: order.Number,
		OrderAmount: order.Amount,
		OrderTime:   order.OrderTime,
	}, nil

}

// RepetitionOrder 再来一单
func (d OrderDao) RepetitionOrder(ctx context.Context, orderId string, userId int) error {
	// 根据订单id查询当前订单详情
	orderDetailList, err := d.GetOrderDetailByOrderId(orderId)
	if err != nil {
		return err
	}
	// 将订单详情对象转换为购物车对象
	var shoppingCartList []model.ShoppingCart
	for _, orderDetail := range orderDetailList {
		// 将原订单详情里面的菜品信息重新复制到购物车对象中
		var shoppingCart model.ShoppingCart
		if err = deepcopier.Copy(&orderDetail).To(&shoppingCart); err != nil {
			return err
		}
		shoppingCart.Id = 0
		shoppingCart.UserId = userId
		shoppingCartList = append(shoppingCartList, shoppingCart)
	}
	// 将购物车对象批量添加到数据库
	if err = d.ShoppingCartDao.InsertBatchShoppingCart(ctx, shoppingCartList); err != nil {
		return err
	}
	return nil
}

// HistoryOrders 分页查询历史订单
func (d OrderDao) HistoryOrders(ctx context.Context, data request.PageQueryOrderDTO, userId int) (*common.PageResult, error) {
	var (
		err        error
		OrderList  []model.Order
		orderVOs   []response.OrderVO
		pageResult common.PageResult
	)
	page, _ := strconv.Atoi(data.Page)
	pageSize, _ := strconv.Atoi(data.PageSize)

	query := d.db.WithContext(ctx).Table("orders")
	if data.Status != "" {
		query = query.Where("status = ?", data.Status)
	}

	if err = query.Count(&pageResult.Total).Error; err != nil {
		return nil, err
	}

	if err = query.Scopes(pageResult.Paginate(&page, &pageSize)).
		Order("order_time desc").
		Scan(&OrderList).Error; err != nil {
		return nil, err
	}

	if OrderList != nil && len(OrderList) > 0 {
		for _, order := range OrderList {
			orderId := strconv.Itoa(order.Id)
			orderDetail, err := d.GetOrderDetailByOrderId(orderId)
			if err != nil {
				return nil, err
			}
			var orderVO response.OrderVO
			if err = deepcopier.Copy(&order).To(&orderVO); err != nil {
				return nil, err
			}
			orderVO.OrderDetailList = orderDetail
			orderVOs = append(orderVOs, orderVO)
		}
	}

	pageResult.Records = orderVOs

	return &pageResult, nil

}

// OrderReminder 用户催单
func (d OrderDao) OrderReminder(orderId string) {
	// 查询订单是否存在
	order, err := d.GetOrderById(orderId)
	if err != nil || order == nil {
		global.Log.Warn("订单错误或者不存在")
	}
	//// 基于WebSocket实现催单
	//jsonMap := map[string]any{
	//	"type":    2,
	//	"orderId": orderId,
	//	"content": "订单号: " + order.Number,
	//}
	//websocket.WSServer.SendToAllClients(jsonMap)
}

// 根据订单id查询订单
func (d OrderDao) GetOrderById(orderId string) (*model.Order, error) {
	var order *model.Order
	if err := global.DB.
		Where("id = ?", orderId).
		Find(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

// UpdateOrder 修改订单信息
func (d OrderDao) UpdateOrder(order *model.Order) error {
	if err := global.DB.Table("orders").
		Updates(order).Error; err != nil {
		return err
	}
	return nil
}

// 根据订单号和用户id查询订单
func (d OrderDao) GetOrderByNumberAndUserId(orderNumber, userId string) (*model.Order, error) {
	var order model.Order
	if err := global.DB.
		Where("number = ?", orderNumber).
		Where("user_id = ?", userId).
		First(&order).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &order, nil
}

// 向订单明细表插入多条数据
func insertBatchOrderDetail(orderDetailList []model.OrderDetail) error {
	if err := global.DB.Create(&orderDetailList).Error; err != nil {
		return err
	}
	return nil
}

// admin

// OrderConfirm 接单
func (d OrderDao) OrderConfirm(ctx context.Context, data request.OrderConfirmDTO) error {
	var id int
	switch data.OrderId.(type) {
	case int:
		id = data.OrderId.(int)
	case string:
		id, _ = strconv.Atoi(data.OrderId.(string))
	case float64:
		id = int(data.OrderId.(float64))
	}
	// 更新订单为已接单
	if err := d.UpdateOrder(&model.Order{
		Id:     id,
		Status: enum.Confirmed,
	}); err != nil {
		return err
	}
	return nil
}

// OrderRejection 拒单
func (d OrderDao) OrderRejection(ctx context.Context, data *request.OrderRejectionDTO) error {
	// 根据id查询订单
	order, err := d.GetOrderById(strconv.Itoa(data.OrderId))
	if err != nil {
		return err
	}
	// 订单只有存在且状态为2（待接单）才可以拒单
	if order == nil || order.Status != enum.ToBeConfirmed {
		return errors.New("错误拒单")
	}
	// 检查支付状态
	if order.PayStatus == enum.Paid {
		// 用户已支付，需要退款
		// 模拟微信退款
		log.Printf("已支付订单被拒单, 退款: [%v￥]", order.Amount)
	}
	// 根据订单id更新订单状态、拒单原因、取消时间
	order.Status = enum.Cancelled
	order.RejectionReason = data.RejectionReason
	order.CancelTime = model.LocalTime(time.Now())
	if err = d.UpdateOrder(order); err != nil {
		return err
	}
	return nil
}

// OrderDelivery 订单派送
func (d OrderDao) OrderDelivery(orderId string) error {
	// 根据id查询订单
	order, _ := d.GetOrderById(orderId)
	// 校验订单是否存在，并且状态为Confirmed
	if order == nil || order.Status != enum.Confirmed {
		return errors.New("订单不存在，或者订单不可派送")
	}
	// 更新订单状态,状态转为派送中
	order.Status = enum.DeliveryInProgress
	if err := d.UpdateOrder(order); err != nil {
		return err
	}
	return nil
}

// OrderComplete 完成订单
func (d OrderDao) OrderComplete(orderId string) error {
	// 根据id查询订单
	order, _ := d.GetOrderById(orderId)
	// 校验订单是否存在，并且状态为DeliveryInProgress
	if order == nil || order.Status != enum.DeliveryInProgress {
		return errors.New("订单不存在，或者订单不可完成")
	}
	// 更新订单状态,状态转为完成
	order.Status = enum.Completed
	order.DeliveryTime = model.LocalTime(time.Now())
	if err := d.UpdateOrder(order); err != nil {
		return err
	}
	return nil
}

// GetDailyTurnover 获取每日营业额
func GetDailyTurnover(begin, end time.Time) (float64, error) {
	var turnover float64
	if err := global.DB.Table("orders").
		Where("status = ?", enum.Completed).
		Where("order_time >= ?", begin).
		Where("order_time <= ?", end).
		Select("ifnull(sum(amount),0) as amount").
		Scan(&turnover).Error; err != nil {
		return 0, err
	}
	return turnover, nil
}

// GetDailyOrderCount 获取每日订单数
func GetDailyOrderCount(begin, end time.Time, status int) (int, error) {
	var orderCount int64
	if status != 0 {
		// 查询有效订单数量
		if err := global.DB.Table("orders").
			Where("order_time >= ?", begin).
			Where("order_time <= ?", end).
			Where("status = ?", status).
			Count(&orderCount).Error; err != nil {
			return 0, err
		}
	} else {
		// 查询总订单数量
		if err := global.DB.Table("orders").
			Where("order_time >= ?", begin).
			Where("order_time <= ?", end).
			Count(&orderCount).Error; err != nil {
			return 0, err
		}
	}
	return int(orderCount), nil
}

// GetSalesTop10 获取销量前十的商品
func GetSalesTop10(begin, end time.Time) ([]request.GoodsSalesDTO, error) {
	goodsSales := make([]request.GoodsSalesDTO, 0)
	if err := global.DB.Table("order_detail").
		Select("order_detail.name, sum(order_detail.number) as number").
		Joins("left join orders on order_detail.order_id = orders.id").
		Where("orders.status = ?", enum.Completed).
		Where("orders.order_time >= ?", begin).
		Where("orders.order_time <= ?", end).
		Group("order_detail.name").
		Order("number desc").
		Limit(10).Offset(0).
		Scan(&goodsSales).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return []request.GoodsSalesDTO{}, err
	}
	return goodsSales, nil
}

// 根据订单详情获取菜品信息字符串
func GetOrderDishes(orderDetails []model.OrderDetail) string {
	orderDishes := ""
	// 将每一条订单菜品信息拼接为字符串（格式：宫保鸡丁*3;）
	for _, orderDetail := range orderDetails {
		orderDishes += orderDetail.Name + "*" + strconv.Itoa(orderDetail.Number) + ";"
	}
	return orderDishes
}

// GetOrderByStatusAndOrderTime 根据状态和下单时间查询订单
func GetOrderByStatusAndOrderTime(status int, orderTime model.LocalTime) ([]model.Order, error) {
	var orders []model.Order
	if err := global.DB.
		Where("status = ?", status).
		Where("order_time < ?", orderTime).
		Find(&orders).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.Order{}, err
	}
	return orders, nil
}

// todo : 催单
