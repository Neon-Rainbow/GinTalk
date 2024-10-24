package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var _ VoteCacheInterface = (*VoteCache)(nil)

type VoteCacheInterface interface {
	UpdatePostHot(ctx context.Context, postID int64, upvote int, createTime time.Time) error
	AddPostHot(ctx context.Context, postID int64, oldUp int, newUp int) error
}

type VoteCache struct {
	*redis.Client
}

func (v *VoteCache) UpdatePostHot(ctx context.Context, postID int64, upvote int, createTime time.Time) error {
	key := GenerateRedisKey(PostRankingTemplate)

	// 使用 Redis Pipeline 更新 ZSet，确保高效和一致性
	pipe := v.TxPipeline()
	pipe.ZAdd(ctx, key, &redis.Z{Score: hot(upvote, createTime), Member: strconv.FormatInt(postID, 10)})

	// 执行 Redis Pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (v *VoteCache) AddPostHot(ctx context.Context, postID int64, oldUp int, newUp int) error {
	key := GenerateRedisKey(PostRankingTemplate)

	// 使用 Redis Pipeline 更新 ZSet，确保高效和一致性
	pipe := v.TxPipeline()
	pipe.ZIncrBy(ctx, key, deltaHot(oldUp, newUp), strconv.FormatInt(postID, 10))

	// 执行 Redis Pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func NewVoteCache(client *redis.Client) VoteCacheInterface {
	return &VoteCache{client}
}
