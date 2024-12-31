package cache

import (
	"GinTalk/DTO"
	"GinTalk/dao/Redis"
	"context"
	"encoding/json"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	_ = iota
	OrderByHot
	OrderByTime

	PostStoreTime = time.Hour * 24 * 7
)

// hot 函数计算一个帖子的热度分数。
//
// 参数：
//   - ups: 赞成票的数量。
//   - date: 帖子的发布时间，使用 time.Time 类型。
//
// 返回值：
//   - 帖子的热度分数，使用 float64 类型。
//
// 计算方法：
//  1. 计算赞成票和反对票的差值。
//  2. 计算票数差值的对数值。
//  3. 根据票数差值的正负确定符号。
//  4. 使用 Unix 时间戳计算时间差值。
//  5. 结合符号、对数值和时间差值计算最终的热度分数。
func hot(ups int, date time.Time) float64 {
	downs := 0
	s := float64(ups - downs)                     // 计算赞成票和反对票的差值
	order := math.Log10(math.Max(math.Abs(s), 1)) // 计算票数的对数

	var sign float64
	if s > 0 {
		sign = 1
	} else if s < 0 {
		sign = -1
	} else {
		sign = 0
	}

	// 使用 Unix 时间戳进行时间计算
	seconds := float64(date.Unix() - 1577808000)

	// 计算热度，并四舍五入到最近的整数
	ans := sign*order + seconds/45000
	return ans
}

// deltaHot 计算两个赞成票数量之间的“热度”差异。
// 它使用赞成票数量的以 10 为底的对数来确定热度的变化。
//
// 参数：
// - oldUp: 之前的赞成票数量。
// - newUp: 新的赞成票数量。
//
// 返回值：
// - 一个表示热度变化的 float64 值。
func deltaHot(oldUp, newUp int) float64 {
	return math.Log10(max(float64(newUp), 1)) - math.Log10(max(float64(oldUp), 1))
}

// SavePost 将帖子存储到 Redis 中
func SavePost(ctx context.Context, summary *DTO.PostSummary) error {
	key := GenerateRedisKey(PostSummaryTemplate, summary.PostID)
	data, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	// 将帖子存储到 Redis 中
	if err := Redis.GetRedisClient().Set(ctx, key, data, PostStoreTime).Err(); err != nil {
		return err
	}

	timestamp := float64(time.Now().Unix())
	hotScore := hot(0, time.Now())

	if err := Redis.GetRedisClient().ZAdd(ctx, GenerateRedisKey(PostTimeTemplate), &redis.Z{
		Score:  timestamp,
		Member: summary.PostID,
	}).Err(); err != nil {
		return err
	}
	if err := Redis.GetRedisClient().ZAdd(ctx, GenerateRedisKey(PostRankingTemplate), &redis.Z{
		Score:  hotScore,
		Member: summary.PostID,
	}).Err(); err != nil {
		return err
	}
	return nil
}

// GetPostIDs 从 Redis 中获取帖子 ID 列表。
// 它使用提供的排序方式、页码和每页帖子数量，从 Redis 中获取帖子 ID 列表。
//
// 参数：
//   - ctx: 操作的上下文，允许取消和超时控制。
//   - order: 帖子的排序方式。其中 1 表示按热度排序，2 表示按时间排序。
//   - pageNum: 分页的页码。
//   - pageSize: 每页的帖子数量。
//
// 返回：
//   - []int64: 一个帖子 ID 的切片。
//   - error: 如果操作失败，则返回错误对象，否则返回 nil。
func GetPostIDs(ctx context.Context, order, pageNum, pageSize int) ([]int64, error) {
	var key string

	caseTemplateMap := map[int]string{
		OrderByHot:  PostRankingTemplate,
		OrderByTime: PostTimeTemplate,
	}

	key = GenerateRedisKey(caseTemplateMap[order])

	// 计算分页的开始和结束位置
	start := int64((pageNum - 1) * pageSize)
	end := start + int64(pageSize) - 1

	// 从 Redis 有序集合中获取帖子 ID 列表
	postIDs, err := Redis.GetRedisClient().ZRevRange(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	resp := make([]int64, len(postIDs))
	for i, id := range postIDs {
		_t, _ := strconv.Atoi(id)
		resp[i] = int64(_t)
	}
	return resp, nil
}

// GetPostSummary 从 Redis 中获取帖子摘要信息。
// 它使用提供的帖子 ID 列表，从 Redis 中获取帖子摘要信息。
//
// 参数:
//   - ctx: 操作的上下文，允许取消和超时控制。
//   - postID: 要获取的帖子 ID 列表。
//
// Returns:
//   - []DTO.PostSummary: 一个帖子摘要信息的切片。
//   - []int64: 无法在 Redis 中找到的帖子 ID 列表。
//   - error: 如果操作失败，则返回错误对象，否则返回 nil。
func GetPostSummary(ctx context.Context, postID []int64) ([]DTO.PostSummary, []int64, error) {
	strKeys := make([]string, len(postID))
	for i, key := range postID {
		strKeys[i] = GenerateRedisKey(PostSummaryTemplate, key)
	}
	values, err := Redis.GetRedisClient().MGet(ctx, strKeys...).Result()
	if err != nil {
		return nil, nil, err
	}
	result := make([]DTO.PostSummary, len(values))
	missingIDs := make([]int64, 0)
	for i, value := range values {
		if value == nil {
			missingIDs = append(missingIDs, postID[i])
			continue
		}
		if err := json.Unmarshal([]byte(value.(string)), &result[i]); err != nil {
			return nil, nil, err
		}
	}
	return result, missingIDs, nil
}

// DeletePost 删除帖子
// 1. 删除帖子的摘要信息
// 2. 删除帖子的时间排序
// 3. 删除帖子的热度排序
func DeletePost(ctx context.Context, postID int64) error {
	key := GenerateRedisKey(PostSummaryTemplate, postID)
	if err := Redis.GetRedisClient().Del(ctx, key).Err(); err != nil {
		return err
	}
	if err := Redis.GetRedisClient().ZRem(ctx, GenerateRedisKey(PostTimeTemplate), postID).Err(); err != nil {
		return err
	}
	if err := Redis.GetRedisClient().ZRem(ctx, GenerateRedisKey(PostRankingTemplate), postID).Err(); err != nil {
		return err
	}
	return nil
}

// DeletePostSummary 删除帖子摘要信息
// 从 Redis 中删除指定帖子的摘要信息。
func DeletePostSummary(ctx context.Context, postID int64) error {
	key := GenerateRedisKey(PostSummaryTemplate, postID)
	if err := Redis.GetRedisClient().Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
