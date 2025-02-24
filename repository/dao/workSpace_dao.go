package dao

import (
	"gorm.io/gorm"
	"takeout/common/enum"
	"takeout/global"
	"takeout/internal/api/admin/response"
	"takeout/internal/model"
	"time"
)

type WorkSpaceDao struct {
	reportDao ReportDao
	db        *gorm.DB
}

func NewWorkSpaceDao() WorkSpaceDao {
	return WorkSpaceDao{
		reportDao: ReportDao{db: global.DB},
		db:        global.DB,
	}
}

// TodayBusinessData 今日数据
func (d *WorkSpaceDao) TodayBusinessData() (response.BusinessDataVO, error) {
	// 今日时间区间
	beginTime := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), 0, 0, 0, 0, time.Local)
	endTime := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), 23, 59, 59, 999999999, time.Local)
	return d.reportDao.GetBusinessData(beginTime, endTime)
}

// OverviewOrders 订单管理
func (d *WorkSpaceDao) OverviewOrders() (response.OrderOverViewVO, error) {
	var (
		waitingOrders   int
		deliveredOrders int
		completedOrders int
		cancelledOrders int
		allOrders       int
		err             error
	)
	// 今日时间区间
	beginTime := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), 0, 0, 0, 0, time.Local)
	endTime := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), 23, 59, 59, 999999999, time.Local)
	// 待接单订单数
	waitingOrders, err = d.reportDao.GetDailyOrderCount(beginTime, endTime, enum.ToBeConfirmed)
	if err != nil {
		return response.OrderOverViewVO{}, err
	}
	// 待派送订单数
	deliveredOrders, err = d.reportDao.GetDailyOrderCount(beginTime, endTime, enum.Confirmed)
	if err != nil {
		return response.OrderOverViewVO{}, err
	}
	// 已完成订单数
	completedOrders, err = d.reportDao.GetDailyOrderCount(beginTime, endTime, enum.Completed)
	if err != nil {
		return response.OrderOverViewVO{}, err
	}
	// 已取消订单数
	cancelledOrders, err = d.reportDao.GetDailyOrderCount(beginTime, endTime, enum.Cancelled)
	if err != nil {
		return response.OrderOverViewVO{}, err
	}
	// 全部订单订单数
	allOrders, err = d.reportDao.GetDailyOrderCount(beginTime, endTime, 0)

	return response.OrderOverViewVO{
		WaitingOrders:   waitingOrders,
		DeliveredOrders: deliveredOrders,
		CompletedOrders: completedOrders,
		CancelledOrders: cancelledOrders,
		AllOrders:       allOrders,
	}, nil
}

// OverviewDishes 菜品总览
func (d *WorkSpaceDao) OverviewDishes() (response.DishOverViewVO, error) {
	sold := GetDishCount(enum.ENABLE)
	discontinued := GetDishCount(enum.DISABLE)
	return response.DishOverViewVO{
		Sold:         sold,
		Discontinued: discontinued,
	}, nil
}

// OverviewSetmeals 套餐总览
func (d *WorkSpaceDao) OverviewSetmeals() (response.SetmealOverViewVO, error) {
	sold := GetSetmealCount(enum.ENABLE)
	discontinued := GetSetmealCount(enum.DISABLE)
	return response.SetmealOverViewVO{
		Sold:         sold,
		Discontinued: discontinued,
	}, nil
}

// GetDailyOrderCount 获取每日订单数
func (d *WorkSpaceDao) GetDailyOrderCount(begin, end time.Time, status int) (int, error) {
	var orderCount int64
	if status != 0 {
		// 查询有效订单数量
		if err := d.db.Table("orders").
			Where("order_time >= ?", model.LocalTime(begin)).
			Where("order_time <= ?", model.LocalTime(end)).
			Where("status = ?", status).
			Count(&orderCount).Error; err != nil {
			return 0, err
		}
	} else {
		// 查询总订单数量
		if err := d.db.Table("orders").
			Where("order_time >= ?", model.LocalTime(begin)).
			Where("order_time <= ?", model.LocalTime(end)).
			Count(&orderCount).Error; err != nil {
			return 0, err
		}
	}
	return int(orderCount), nil
}

// GetDishCount 获取某状态菜品数量
func GetDishCount(status int) int {
	var dishCount int64
	if err := global.DB.Table("dish").
		Where("status = ?", status).
		Count(&dishCount).Error; err != nil {
		return 0
	}
	return int(dishCount)
}

// GetSetmealCount 获取某状态套餐数量
func GetSetmealCount(status int) int {
	var setmealCount int64
	if err := global.DB.Table("setmeal").
		Where("status = ?", status).
		Count(&setmealCount).Error; err != nil {
		return 0
	}
	return int(setmealCount)
}
