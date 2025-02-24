package repository

import (
	"takeout/internal/api/admin/response"
	"time"
)

type WorkSpaceRepo interface {
	OverviewSetmeals() (response.SetmealOverViewVO, error)
	GetDailyOrderCount(begin, end time.Time, status int) (int, error)
	TodayBusinessData() (response.BusinessDataVO, error)
	OverviewOrders() (response.OrderOverViewVO, error)
	OverviewDishes() (response.DishOverViewVO, error)
}
