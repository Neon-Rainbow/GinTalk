package cache

import (
	"GinTalk/DTO"
	"GinTalk/dao/Redis"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"math"
	"strconv"
	"time"
)

const (
	NoOrder = iota
	OrderByHot
	OrderByTime

	PostStoreTime = time.Hour * 24 * 7
)

// hot 用于计算帖子的热度
// Reddit 热度算法实现
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

func DeletePostSummary(ctx context.Context, postID int64) error {
	key := GenerateRedisKey(PostSummaryTemplate, postID)
	if err := Redis.GetRedisClient().Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
