package response

// BusinessDataVO 工作台今日数据概览
type BusinessDataVO struct {
	Turnover            float64 `json:"turnover"`
	ValidOrderCount     int     `json:"validOrderCount"`
	OrderCompletionRate float64 `json:"orderCompletionRate"`
	UnitPrice           float64 `json:"unitPrice"`
	NewUsers            int     `json:"newUsers"`
}

// OrderOverViewVO 订单概览数据
type OrderOverViewVO struct {
	WaitingOrders   int `json:"waitingOrders"`
	DeliveredOrders int `json:"deliveredOrders"`
	CompletedOrders int `json:"completedOrders"`
	CancelledOrders int `json:"cancelledOrders"`
	AllOrders       int `json:"allOrders"`
}

// DishOverViewVO 菜品数据总览
type DishOverViewVO struct {
	Sold         int `json:"sold"`
	Discontinued int `json:"discontinued"`
}

// SetmealOverViewVO 套餐数据总览
type SetmealOverViewVO struct {
	Sold         int `json:"sold"`
	Discontinued int `json:"discontinued"`
}
