package admin

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type EmployeeRouter struct {
	service service.IEmployeeService
}

func (er *EmployeeRouter) InitApiRouter(router *gin.RouterGroup) {
	// admin/employee
	publicRouter := router.Group("employee")
	privateRouter := router.Group("employee")
	// 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTAdmin())
	// 依赖注入
	er.service = service.NewEmployeeService(dao.NewEmployeeDao(global.DB))
	employeeCtl := controller.NewEmployeeController(er.service)
	{
		publicRouter.POST("/login", employeeCtl.Login)
		privateRouter.POST("/logout", employeeCtl.Logout)
		privateRouter.POST("", employeeCtl.AddEmployee)

		privateRouter.GET("/page", employeeCtl.PageQuery)

		privateRouter.POST("/status/:status", employeeCtl.OnOrOff)
		privateRouter.PUT("/editPassword", employeeCtl.EditPassword)
		privateRouter.PUT("", employeeCtl.UpdateEmployee)
		privateRouter.GET("/:id", employeeCtl.GetById)
	}
}
