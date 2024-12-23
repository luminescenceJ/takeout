package initialize

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "takeout/docs" // 导入生成的 Swagger 文档
	"takeout/internal/router"
)

func routerInit() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	allRouter := router.AllRouter
	// admin
	admin := r.Group("/admin")
	{
		//employee 路由
		allRouter.EmployeeRouter.InitApiRouter(admin)
	}
	return r
}
