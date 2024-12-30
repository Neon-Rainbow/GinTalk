package cache

import (
	"GinTalk/dao/Redis"
	"context"
	"time"
)

// AddTokenToBlacklist 将token添加到Redis黑名单中，并指定过期时间。
// 它使用提供的token和预定义模板生成Redis键，然后将此键设置为值"1"并指定过期时间。
//
// 参数:
//   - ctx: 操作的上下文，允许取消和超时控制。
//   - token: 要加入黑名单的token。
//   - expiration: token在黑名单中应保留的时间。
//
// 返回:
//   - error: 如果操作失败，则返回错误对象，否则返回nil。
func AddTokenToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := GenerateRedisKey(BlackListTokenKeyTemplate, token)
	err := Redis.GetRedisClient().Set(ctx, key, "1", expiration).Err()
	return err
}

// IsTokenInBlacklist 判断token是否在黑名单中。
// 它使用提供的token生成Redis键，并检查其在Redis中的存在性。
//
// 参数:
//   - ctx: 操作的上下文，允许取消和超时控制。
//   - token: 要检查的token。
//
// 返回:
//   - bool: 如果token在黑名单中，则返回true，否则返回false。
//   - error: 如果操作失败，则返回错误对象，否则返回nil。
//
// 使用示例:
//
//	isInBlacklist, err := cache.IsTokenInBlacklist(ctx, token)
//	if err != nil {
//		// 处理错误
//	}
//	if isInBlacklist {
//		// token在黑名单中
//	}
func IsTokenInBlacklist(ctx context.Context, token string) (bool, error) {
	key := GenerateRedisKey(BlackListTokenKeyTemplate, token)
	exists, err := Redis.GetRedisClient().Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
