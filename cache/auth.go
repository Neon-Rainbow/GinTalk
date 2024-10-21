package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type AuthCacheInterface interface {
	// SetUserToken 设置用户token
	// userID: 用户ID
	// accessToken: 访问令牌
	// refreshToken: 刷新令牌
	// 如果设置成功, 返回nil
	// 由于项目使用 Token 黑名单机制,因此该接口以已经废弃
	SetUserToken(ctx context.Context, userID int64, accessToken, refreshToken string) error

	// AddTokenToBlacklist 将token加入黑名单
	// expiration: 过期时间
	// token: token
	// 如果加入成功, 返回nil
	AddTokenToBlacklist(ctx context.Context, token string, expiration time.Duration) error

	// IsTokenInBlacklist 判断token是否在黑名单中
	// 如果在黑名单中, 返回true
	// 如果不在黑名单中, 返回false
	IsTokenInBlacklist(ctx context.Context, token string) (bool, error)
}

type AuthCache struct {
	*redis.Client
}

// NewAuthCache 实例化 AuthCache
func NewAuthCache(client *redis.Client) AuthCacheInterface {
	return &AuthCache{client}
}

func (a *AuthCache) SetUserToken(ctx context.Context, userId int64, accessToken string, refreshToken string) error {
	//key := GenerateRedisKey(UserTokenKeyTemplate, userId)
	//err := a.HSet(ctx, key, map[string]interface{}{
	//	"accessToken":  accessToken,
	//	"refreshToken": refreshToken,
	//}).Err()
	//return err
	return fmt.Errorf(" SetUserToken(ctx context.Context, userId int64, accessToken string, refreshToken string) error 该接口已经废弃,请使用其他接口")
}

func (a *AuthCache) AddTokenToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := GenerateRedisKey(BlackListTokenKeyTemplate, token)
	err := a.Set(ctx, key, "1", expiration).Err()
	return err
}

func (a *AuthCache) IsTokenInBlacklist(ctx context.Context, token string) (bool, error) {
	key := GenerateRedisKey(BlackListTokenKeyTemplate, token)
	return a.Exists(ctx, key).Val() == 1, nil
}
