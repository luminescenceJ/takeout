package repository

import (
	"takeout/internal/api/admin/response"
	"takeout/internal/api/user/request"
	"time"
)

type ReportRepo interface {
	GetDailyTurnover(begin, end time.Time) (float64, error)
	GetDailyOrderCount(begin, end time.Time, status int) (int, error)
	GetSalesTop10(begin, end time.Time) ([]request.GoodsSalesDTO, error)
	GetUserCount(begin, end time.Time) (int, error)
	GetBusinessData(beginTime, endTime time.Time) (response.BusinessDataVO, error)
}
