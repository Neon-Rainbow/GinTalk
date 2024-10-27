package cache

import (
	"GinTalk/dao/Redis"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

// UpdatePostHot 更新帖子的热度
func UpdatePostHot(ctx context.Context, postID int64, upvote int, createTime time.Time) error {
	key := GenerateRedisKey(PostRankingTemplate)

	// 使用 Redis Pipeline 更新 ZSet，确保高效和一致性
	pipe := Redis.GetRedisClient().TxPipeline()
	pipe.ZAdd(ctx, key, &redis.Z{Score: hot(upvote, createTime), Member: strconv.FormatInt(postID, 10)})

	// 执行 Redis Pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

// AddPostHot 增加帖子的热度
func AddPostHot(ctx context.Context, postID int64, oldUp int, newUp int) error {
	key := GenerateRedisKey(PostRankingTemplate)

	// 使用 Redis Pipeline 更新 ZSet，确保高效和一致性
	pipe := Redis.GetRedisClient().TxPipeline()
	pipe.ZIncrBy(ctx, key, deltaHot(oldUp, newUp), strconv.FormatInt(postID, 10))

	// 执行 Redis Pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}
