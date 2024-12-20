package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/controller"
	"takeout/internal/service"
	"takeout/repository/dao"
)

type EmployeeRouter struct {
	service service.IEmployeeService
}

func (er *EmployeeRouter) InitApiRouter(router *gin.RouterGroup) {
	// admin/employee
	publicRouter := router.Group("employee")
	//privateRouter := router.Group("employee")
	er.service = service.NewEmployeeService(dao.NewEmployeeDao(global.DB))
	employeeCtl := controller.NewEmployeeController(er.service)
	{
		publicRouter.POST("/login", employeeCtl.Login)
	}
}
