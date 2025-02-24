package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/admin/response"
	"takeout/repository/dao"
)

type WorkspaceController struct {
	dao dao.WorkSpaceDao
}

func NewWorkSpaceController(dao dao.WorkSpaceDao) *WorkspaceController {
	return &WorkspaceController{
		dao: dao,
	}
}

// TodayBusinessData 今日数据
func (c *WorkspaceController) TodayBusinessData(ctx *gin.Context) {

	var (
		businessData response.BusinessDataVO
		code         = e.SUCCESS
		err          error
	)

	// 调用service层进行处理
	if businessData, err = c.dao.TodayBusinessData(); err != nil {
		code = e.ERROR
		global.Log.Warn("TodayBusinessData Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
		return
	}
	// 日志打印
	log.Printf("今日数据: %v", businessData)

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: businessData,
		Msg:  e.GetMsg(code),
	})
}

// OverviewOrders 订单管理
func (c *WorkspaceController) OverviewOrders(ctx *gin.Context) {
	var (
		ordersData response.OrderOverViewVO
		code       = e.SUCCESS
		err        error
	)

	// 调用service层进行处理
	ordersData, err = c.dao.OverviewOrders()
	if err != nil {
		code = e.ERROR
		global.Log.Warn("OverviewOrders Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
	}
	// 日志打印
	log.Printf("订单管理: %v", ordersData)

	// 响应
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: ordersData,
		Msg:  e.GetMsg(code),
	})
}

// OverviewDishes 菜品总览
func (c *WorkspaceController) OverviewDishes(ctx *gin.Context) {
	var (
		dishesData response.DishOverViewVO
		code       = e.SUCCESS
		err        error
	)

	// 调用service层进行处理
	dishesData, err = c.dao.OverviewDishes()
	if err != nil {
		code = e.ERROR
		global.Log.Warn("OverviewDishes Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
	}
	// 日志打印
	log.Printf("菜品总览: %v", dishesData)

	// 响应
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: dishesData,
		Msg:  e.GetMsg(code),
	})

}

// OverviewSetmeals 套餐总览
func (c *WorkspaceController) OverviewSetmeals(ctx *gin.Context) {

	var (
		setmealsData response.SetmealOverViewVO
		code         = e.SUCCESS
		err          error
	)

	// 调用service层进行处理
	setmealsData, err = c.dao.OverviewSetmeals()
	if err != nil {
		code = e.ERROR
		global.Log.Warn("OverviewSetmeals Error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  e.GetMsg(code),
		})
	}
	// 日志打印
	log.Printf("套餐总览: %v", setmealsData)

	// 响应
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: setmealsData,
		Msg:  e.GetMsg(code),
	})
}
