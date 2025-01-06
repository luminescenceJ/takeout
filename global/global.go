package global

import (
	"gorm.io/gorm"
	"takeout/config"
	"takeout/logger"
)

var (
	DB     *gorm.DB
	Log    logger.ILog
	Config *config.AllConfig
	Path   *config.Path
)
