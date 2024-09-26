package dao

import (
	"context"
	"fmt"
	"GinTalk/dao/Redis"
	"time"
)

// StoreKeyValue 用于存储键值对
func StoreKeyValue[K uint | string, V string | int | int64 | float64 | uint | map[string]interface{} | []string](ctx context.Context, key K, value V, expiration time.Duration) error {
	rdb := Redis.GetRedisClient()
	stringKey := string(key)
	switch v := any(value).(type) {
	case string:
		// 存储字符串类型
		err := rdb.Set(ctx, stringKey, v, expiration).Err()
		if err != nil {
			return err
		}
	case int, int64, float64:
		// 存储数字类型，转换为字符串
		err := rdb.Set(ctx, stringKey, fmt.Sprintf("%v", v), expiration).Err()
		if err != nil {
			return err
		}
	case map[string]interface{}:
		// 存储哈希类型
		err := rdb.HMSet(ctx, stringKey, v).Err()
		if err != nil {
			return err
		}
	case []string:
		// 存储列表类型
		err := rdb.RPush(ctx, stringKey, v).Err()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported value type: %V", v)
	}
	return nil
}

// GetKeyValue 用于获取键值对
func GetKeyValue[T any](ctx context.Context, key string) (T, error) {
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
