package Redis

import (
	"context"
	"fmt"
	"GinTalk/settings"
	"github.com/go-redis/redis/v8"
)

// redisClient 用于存储redis连接
var redisClient *redis.Client

// Init 初始化redis连接
func Init(config *settings.RedisConfig) (err error) {
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
	return nil
}

// Close 关闭redis连接
func Close() {
	_ = redisClient.Close()
}

func GetRedisClient() *redis.Client {
	return redisClient
}
