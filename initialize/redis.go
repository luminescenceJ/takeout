package initialize

import (
	"errors"
	"github.com/go-redis/redis"
	"takeout/global"
)

func InitRedis() *redis.Client {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Host + ":" + global.Config.Redis.Port,
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.Database,
	})
	if _, err := RedisClient.Ping().Result(); err != nil {
		panic(errors.New("init redis fail:" + err.Error()))
	}
	return RedisClient
}
