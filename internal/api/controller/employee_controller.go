package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/request"
	"takeout/internal/service"
)

type EmployeeController struct {
	service service.IEmployeeService
}

func NewEmployeeController(employeeService service.IEmployeeService) *EmployeeController {
	return &EmployeeController{service: employeeService}
}

func (ec *EmployeeController) Login(ctx *gin.Context) {
	code := e.SUCCESS
	employeeLogin := request.EmployeeLogin{}
	err := ctx.ShouldBindWith(&employeeLogin, binding.JSON)
	if err != nil {
		code = e.ERROR
		global.Log.Debug("EmployeeController login 解析失败")
		return
	}

	resp, err := ec.service.Login(ctx, employeeLogin)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("EmployeeController login Error:", err.Error())
		ctx.JSON(http.StatusOK, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: resp,
	})
}
