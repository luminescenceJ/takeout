package utils

import (
	"context"
	"takeout/global"
)

type RedisBloom struct {
	Key string // Redis中的布隆过滤器键名
	Ctx context.Context
}

// 创建新的 RedisBloom 实例（不重新建立连接）
func NewRedisBloomFilter(key string) *RedisBloom {
	return &RedisBloom{
		Key: key,
		Ctx: context.Background(),
	}
}

// 初始化布隆过滤器
func (r *RedisBloom) Init(errRate float64, capacity int64) error {
	_, err := global.RedisClient.Do(r.Ctx, "BF.RESERVE", r.Key, errRate, capacity).Result()
	return err // 如果 key 已存在，这里会报错，可以忽略处理
}

// 添加菜品ID
func (r *RedisBloom) AddDish(dishID string) error {
	_, err := global.RedisClient.Do(r.Ctx, "BF.ADD", r.Key, dishID).Result()
	return err
}

// 查询菜品ID是否可能存在
func (r *RedisBloom) HasDish(dishID string) (bool, error) {
	return global.RedisClient.Do(r.Ctx, "BF.EXISTS", r.Key, dishID).Bool()
}
