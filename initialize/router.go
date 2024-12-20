package initialize

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/router"
)

func routerInit() *gin.Engine {
	r := gin.Default()
	allRouter := router.AllRouter
	// admin
	admin := r.Group("/admin")
	{
		//employee 路由
		allRouter.EmployeeRouter.InitApiRouter(admin)
	}
	return r
}
