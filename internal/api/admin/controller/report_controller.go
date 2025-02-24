package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/admin/response"
	"takeout/internal/service"
)

type ReportController struct {
	service service.IReportService
}

func NewReportController(service service.IReportService) *ReportController {
	return &ReportController{service: service}
}

// TurnoverStatistics 营业额数据统计
func (c *ReportController) TurnoverStatistics(ctx *gin.Context) {
	var (
		turnover response.TurnoverReportVO
		code     = e.SUCCESS
		err      error
	)

	// 接收query参数
	begin := ctx.Query("begin")
	end := ctx.Query("end")

	// 调用service层进行处理
	if turnover, err = c.service.TurnoverStatistics(begin, end); err != nil {
		code = e.ERROR
		global.Log.Warn("TurnoverStatistics Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	// 日志打印
	log.Printf("营业额数据统计: %v", turnover)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: turnover,
		Msg:  e.GetMsg(code),
	})
}

// UserStatistics 用户统计
func (c *ReportController) UserStatistics(ctx *gin.Context) {
	var (
		userStatistics response.UserReportVO
		code           = e.SUCCESS
		err            error
	)

	// 接收query参数
	begin := ctx.Query("begin")
	end := ctx.Query("end")

	// 调用service层进行处理
	if userStatistics, err = c.service.UserStatistics(begin, end); err != nil {
		code = e.ERROR
		global.Log.Warn("UserStatistics Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	// 日志打印
	log.Printf("用户统计: %v", userStatistics)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: userStatistics,
		Msg:  e.GetMsg(code),
	})
}

// ReportOrderStatistics 订单统计
func (c *ReportController) ReportOrderStatistics(ctx *gin.Context) {
	var (
		data response.OrderReportVO
		code = e.SUCCESS
		err  error
	)

	// 接收query参数
	begin := ctx.Query("begin")
	end := ctx.Query("end")

	// 调用service层进行处理
	if data, err = c.service.ReportOrderStatistics(begin, end); err != nil {
		code = e.ERROR
		global.Log.Warn("ReportOrderStatistics Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	// 日志打印
	log.Printf("订单统计: %v", data)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: data,
		Msg:  e.GetMsg(code),
	})
}

// Top10Statistics 销量排名
func (c *ReportController) Top10Statistics(ctx *gin.Context) {
	var (
		data response.SalesTop10ReportVO
		code = e.SUCCESS
		err  error
	)

	// 接收query参数
	begin := ctx.Query("begin")
	end := ctx.Query("end")

	// 调用service层进行处理
	if data, err = c.service.Top10Statistics(begin, end); err != nil {
		code = e.ERROR
		global.Log.Warn("Top10Statistics Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	// 日志打印
	log.Printf("销量排名: %v", data)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: data,
		Msg:  e.GetMsg(code),
	})

}

// ExportExcel 导出运营数据Excel报表
func (c *ReportController) ExportExcel(ctx *gin.Context) {
	// 调用service层进行处理
	c.service.ExportExcel(ctx)
	// 日志打印
	log.Printf("导出运营数据Excel报表")
}
