package cache

import (
	"GinTalk/dao/Redis"
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// UpdatePostHot 函数用于更新 Redis 有序集合中帖子的热度分数。
//
// 该函数使用 Redis 管道来确保高效和一致的 ZSet 更新。热度分数是根据点赞数和
// 帖子的创建时间计算的。
//
// 参数:
//   - ctx: 请求的上下文，用于取消和超时控制。
//   - postID: 帖子的唯一标识符。
//   - upvote: 帖子收到的点赞数。
//   - createTime: 帖子的创建时间。
//
// 返回值:
//   - error: 如果 Redis 管道执行失败，则返回错误，否则返回 nil。
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

// AddPostHot 函数使用 Redis 管道来确保高效和一致的 ZSet 更新。
//
// 参数:
//   - ctx: 请求的上下文，用于取消和超时控制。
//   - postID: 帖子的唯一标识符。
//   - oldUp: 帖子之前的点赞数。
//   - newUp: 帖子新的点赞数。
//
// 返回值:
//   - error: 如果 Redis 管道执行失败，则返回错误，否则返回 nil。
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
