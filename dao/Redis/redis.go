package Redis

import (
	"GinTalk/settings"
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// redisClient 用于存储redis连接
var redisClient *redis.Client
var once sync.Once

// initRedis 初始化redis连接
func initRedis(config *settings.RedisConfig) (err error) {
	// 初始化redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: "",
		DB:       0,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	redisClient = rdb

	// 向redis中写入数据
	err = rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Close 关闭redis连接
func Close() {
	err := redisClient.Close()
	if err != nil {
		zap.L().Error("关闭redis连接失败", zap.Error(err))
	}
}

// GetRedisClient 获取redis连接
func GetRedisClient() *redis.Client {
	once.Do(
		func() {
			err := initRedis(settings.GetConfig().RedisConfig)
			if err != nil {
				zap.L().Fatal("初始化redis连接失败", zap.Error(err))
			}
		})
	return redisClient
}
