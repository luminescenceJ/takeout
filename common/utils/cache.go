package utils

import "takeout/global"

func CleanCache(pattern string) {
	keys, _ := global.RedisClient.Keys(pattern).Result()
	global.RedisClient.Del(keys...)
}
