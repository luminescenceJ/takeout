package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/admin/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type ReportRouter struct {
	service service.IReportService
}

func (er *ReportRouter) InitApiRouter(router *gin.RouterGroup) {
	// /admin/report
	privateRouter := router.Group("report")

	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTAdmin())

	// 依赖注入
	er.service = service.NewReportService(dao.NewReportDao(global.DB))
	reportCtl := controller.NewReportController(er.service)
	{
		// 营业额数据统计
		privateRouter.GET("turnoverStatistics", reportCtl.TurnoverStatistics)
		// 用户统计
		privateRouter.GET("userStatistics", reportCtl.UserStatistics)
		// 订单统计
		privateRouter.GET("ordersStatistics", reportCtl.ReportOrderStatistics)
		// 销量排名
		privateRouter.GET("top10", reportCtl.Top10Statistics)
		// 导出运营数据Excel报表
		privateRouter.GET("export", reportCtl.ExportExcel)
	}
}
