package utils

import "github.com/go-redis/redis/v8"

// Formattable 用于限定可以格式化的类型
type Formattable interface {
	~int | ~int32 | ~int64 | ~uint | ~uint64 | string
}

// RedisValueFormattable 用于限定可以存入 redis 的数据类型
type RedisValueFormattable interface {
	string | int | int64 | float64 | uint | map[string]interface{} | []string | []*redis.Z
}
