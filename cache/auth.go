package cache

import (
	"GinTalk/dao/Redis"
	"context"
	"time"
)

// AddTokenToBlacklist 将token加入黑名单
func AddTokenToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := GenerateRedisKey(BlackListTokenKeyTemplate, token)
	err := Redis.GetRedisClient().Set(ctx, key, "1", expiration).Err()
	return err
}

// IsTokenInBlacklist 判断token是否在黑名单中
func IsTokenInBlacklist(ctx context.Context, token string) (bool, error) {
	key := GenerateRedisKey(BlackListTokenKeyTemplate, token)
	return Redis.GetRedisClient().Exists(ctx, key).Val() == 1, nil
}
