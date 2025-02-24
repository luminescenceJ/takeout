package initialize

import (
	"github.com/gin-gonic/gin"
	"takeout/config"
	"takeout/global"
	"takeout/logger"
)

func GlobalInit() *gin.Engine {

	// 配置文件初始化
	global.Config = config.InitLoadConfig()
	// Log初始化
	global.Log = logger.NewLogger(global.Config.Log.Level, global.Config.Log.FilePath)
	// Gorm初始化
	global.DB = InitDatabase(global.Config.DataSource.Dsn())
	// Redis初始化
	global.RedisClient = InitRedis()
	// Router初始化
	router := routerInit()
	return router
}
