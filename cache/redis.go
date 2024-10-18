package cache

import (
	"GinTalk/dao/Redis"
	"GinTalk/utils"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// storeKeyValue 用于存储键值对，支持 string、int、float、hash、list、zset
func storeKeyValue[KeyTemplate utils.Formattable, ValueTemplate utils.RedisValueFormattable](
	ctx context.Context,
	key KeyTemplate,
	value ValueTemplate,
	expiration time.Duration,
) error {
	rdb := Redis.GetRedisClient()
	stringKey := fmt.Sprintf("%v", key)

	switch v := any(value).(type) {
	case string:
		// 存储字符串类型
		err := rdb.Set(ctx, stringKey, v, expiration).Err()
		if err != nil {
			return err
		}
	case int, int64, float64, uint:
		// 存储数字类型，转换为字符串
		err := rdb.Set(ctx, stringKey, fmt.Sprintf("%v", v), expiration).Err()
		if err != nil {
			return err
		}
	case map[string]interface{}:
		// 存储哈希类型
		err := rdb.HSet(ctx, stringKey, v).Err()
		if err != nil {
			return err
		}
	case []string:
		// 存储列表类型
		err := rdb.RPush(ctx, stringKey, v).Err()
		if err != nil {
			return err
		}
	case []*redis.Z:
		// 存储有序集合 ZSET 类型
		err := rdb.ZAdd(ctx, stringKey, v...).Err() // 解包切片为可变参数
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported value type: %T", v)
	}
	return nil
}

// getKeyValue 用于获取键值对
func getKeyValue[T utils.RedisValueFormattable](ctx context.Context, key string) (T, error) {
	rdb := Redis.GetRedisClient()
	var result T
	switch any(result).(type) {
	case string:
		// 获取字符串类型数据
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return result, err
		}
		result = any(val).(T)
	case int, int64, float64:
		// 获取数字类型数据（Redis 中存储为字符串）
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return result, err
		}
		var numVal T
		switch any(numVal).(type) {
		case int:
			var parsedVal int
			fmt.Sscanf(val, "%d", &parsedVal)
			result = any(parsedVal).(T)
		case int64:
			var parsedVal int64
			fmt.Sscanf(val, "%d", &parsedVal)
			result = any(parsedVal).(T)
		case float64:
			var parsedVal float64
			fmt.Sscanf(val, "%f", &parsedVal)
			result = any(parsedVal).(T)
		}
	case map[string]string:
		// 获取哈希类型数据
		val, err := rdb.HGetAll(ctx, key).Result()
		if err != nil {
			return result, err
		}
		result = any(val).(T)
	case []string:
		// 获取列表类型数据
		val, err := rdb.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			return result, err
		}
		result = any(val).(T)
	default:
		return result, fmt.Errorf("unsupported value type: %T", result)
	}
	return result, nil
}

// Invalidate 使缓存失效
func Invalidate(ctx context.Context, key string) error {
	err := Redis.GetRedisClient().Del(ctx, key).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}

// InvalidatePattern 删除 Redis 中符合模式的缓存数据
// @Param ctx context.Context 上下文
// @Param pattern string 模式
// @Return error 错误信息
// 例如：InvalidatePattern(ctx, "user:*") 删除所有以 user: 开头的键
func InvalidatePattern(ctx context.Context, pattern string) error {
	iter := Redis.GetRedisClient().Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := Redis.GetRedisClient().Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

// InvalidateAll 删除 Redis 中所有缓存数据
func InvalidateAll(ctx context.Context) error {
	err := Redis.GetRedisClient().FlushDB(ctx).Err()
	if err != nil {
		return err
	}
	return nil
}
