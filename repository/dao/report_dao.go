package dao

import (
	"gorm.io/gorm"
	"takeout/common/enum"
	"takeout/internal/api/admin/response"
	"takeout/internal/api/user/request"
	"takeout/internal/model"
	"time"
)

type ReportDao struct {
	db *gorm.DB
}

func NewReportDao(db *gorm.DB) *ReportDao {
	return &ReportDao{db: db}
}

// GetDailyTurnover 获取每日营业额
func (d *ReportDao) GetDailyTurnover(begin, end time.Time) (float64, error) {
	var turnover float64
	if err := d.db.Table("orders").
		Where("status = ?", enum.Completed).
		Where("order_time >= ?", model.LocalTime(begin)).
		Where("order_time <= ?", model.LocalTime(end)).
		Select("ifnull(sum(amount),0) as amount").
		Scan(&turnover).Error; err != nil {
		return 0, err
	}
	return turnover, nil
}

// GetDailyOrderCount 获取每日订单数
func (d *ReportDao) GetDailyOrderCount(begin, end time.Time, status int) (int, error) {
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

// GetSalesTop10 获取销量前十的商品
func (d *ReportDao) GetSalesTop10(begin, end time.Time) ([]request.GoodsSalesDTO, error) {
	goodsSales := make([]request.GoodsSalesDTO, 0)
	if err := d.db.Table("order_detail").
		Select("order_detail.name, sum(order_detail.number) as number").
		Joins("left join orders on order_detail.order_id = orders.id").
		Where("orders.status = ?", enum.Completed).
		Where("orders.order_time >= ?", begin).
		Where("orders.order_time <= ?", end).
		Group("order_detail.name").
		Order("number desc").
		Limit(10).Offset(0).
		Scan(&goodsSales).Error; err != nil {
		return []request.GoodsSalesDTO{}, err
	}
	return goodsSales, nil
}

// GetUserCount 获取用户数量
func (d *ReportDao) GetUserCount(begin, end time.Time) (int, error) {
	var userCount int64
	if begin.IsZero() {
		// 查询总用户数量
		if err := d.db.Table("user").
			Where("create_time <= ?", model.LocalTime(end)).
			Count(&userCount).Error; err != nil {
			return 0, err
		}
	} else {
		// 查询新增用户数量
		if err := d.db.Table("user").
			Where("create_time >= ?", model.LocalTime(begin)).
			Where("create_time <= ?", model.LocalTime(end)).
			Count(&userCount).Error; err != nil {
			return 0, err
		}
	}
	return int(userCount), nil
}

// GetBusinessData 获取营业数据
func (d *ReportDao) GetBusinessData(beginTime, endTime time.Time) (response.BusinessDataVO, error) {
	// 查询总订单数
	totalOrderCount, _ := d.GetDailyOrderCount(beginTime, endTime, 0)
	// 营业额
	turnover, _ := d.GetDailyTurnover(beginTime, endTime)
	// 有效订单数
	validOrderCount, _ := d.GetDailyOrderCount(beginTime, endTime, enum.Completed)
	// 订单完成率
	var orderCompletionRate float64
	// 平均客单价
	var unitPrice float64
	if totalOrderCount != 0 && validOrderCount != 0 {
		orderCompletionRate = float64(validOrderCount) / float64(totalOrderCount)
		unitPrice = turnover / float64(validOrderCount)
	}
	// 新增用户数
	newUsers, _ := d.GetUserCount(beginTime, endTime)
	return response.BusinessDataVO{
		Turnover:            turnover,
		ValidOrderCount:     validOrderCount,
		OrderCompletionRate: orderCompletionRate,
		UnitPrice:           unitPrice,
		NewUsers:            newUsers,
	}, nil
}
