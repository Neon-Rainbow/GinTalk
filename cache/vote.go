package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var _ VoteCacheInterface = (*VoteCache)(nil)

type VoteCacheInterface interface {
	UpdatePostHot(ctx context.Context, postID int64, upvote int, downvote int, createTime time.Time) error
}

type VoteCache struct {
	*redis.Client
}

func (v *VoteCache) UpdatePostHot(ctx context.Context, postID int64, upvote int, downvote int, createTime time.Time) error {
	key := GenerateRedisKey(PostRankingTemplate)

	// 重新计算帖子的热度
	ranking := hot(upvote, downvote, createTime)

	// 使用 Redis Pipeline 更新 ZSet，确保高效和一致性
	pipe := v.TxPipeline()
	pipe.ZAdd(ctx, key, &redis.Z{
		Score:  ranking,
		Member: postID,
	})

	// 执行 Redis Pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to update post hot score: %w", err)
	}
	return nil
}

func NewVoteCache(client *redis.Client) VoteCacheInterface {
	return &VoteCache{client}
}
