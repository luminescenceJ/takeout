package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"takeout/config"
	"takeout/logger"
)

var (
	DB          *gorm.DB
	Log         logger.ILog
	Config      *config.AllConfig
	Path        *config.Path
	RedisClient *redis.Client
)
