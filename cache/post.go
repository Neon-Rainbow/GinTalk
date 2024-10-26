package cache

import (
	"GinTalk/DTO"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"math"
	"strconv"
	"time"
)

var _ PostCacheInterface = (*PostCache)(nil)

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

type PostCacheInterface interface {
	SavePostToRedis(ctx context.Context, summary *DTO.PostSummary) error
	GetPostIDsFromRedis(ctx context.Context, order, pageNum, pageSize int) ([]int64, error)
	GetPostSummaryFromRedis(ctx context.Context, postID []int64) (postList []DTO.PostSummary, missingIDs []int64, err error)

	// DeleteRedisPost 删除 Redis 中的帖子
	// 该方法会删除帖子的摘要信息、帖子的发布时间和帖子的热度信息
	// 但是该方法不会删除帖子的对应的评论信息
	DeleteRedisPost(ctx context.Context, postID int64) error

	// DeleteRedisPostSummary 删除 Redis 中的帖子摘要信息
	// 该方法会删除帖子的摘要信息，但是不会删除帖子的发布时间和帖子的热度信息
	// 也不会删除帖子的对应的评论信息
	// 该接口用于在帖子的摘要信息发生变化时，删除 Redis 中的旧摘要信息
	DeleteRedisPostSummary(ctx context.Context, postID int64) error
}

type PostCache struct {
	*redis.Client
}

func (pc *PostCache) SavePostToRedis(ctx context.Context, summary *DTO.PostSummary) error {
	key := GenerateRedisKey(PostSummaryTemplate, summary.PostID)
	data, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	// 将帖子存储到 Redis 中
	if err := pc.Set(ctx, key, data, PostStoreTime).Err(); err != nil {
		return err
	}

	timestamp := float64(time.Now().Unix())
	hotScore := hot(0, time.Now())

	if err := pc.ZAdd(ctx, GenerateRedisKey(PostTimeTemplate), &redis.Z{
		Score:  timestamp,
		Member: summary.PostID,
	}).Err(); err != nil {
		return err
	}
	if err := pc.ZAdd(ctx, GenerateRedisKey(PostRankingTemplate), &redis.Z{
		Score:  hotScore,
		Member: summary.PostID,
	}).Err(); err != nil {
		return err
	}
	return nil
}

func (pc *PostCache) GetPostIDsFromRedis(ctx context.Context, order, pageNum, pageSize int) ([]int64, error) {
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
	postIDs, err := pc.ZRevRange(ctx, key, start, end).Result()
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

func (pc *PostCache) GetPostSummaryFromRedis(ctx context.Context, postID []int64) ([]DTO.PostSummary, []int64, error) {
	strKeys := make([]string, len(postID))
	for i, key := range postID {
		strKeys[i] = GenerateRedisKey(PostSummaryTemplate, key)
	}
	values, err := pc.MGet(ctx, strKeys...).Result()
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

func (pc *PostCache) DeleteRedisPost(ctx context.Context, postID int64) error {
	key := GenerateRedisKey(PostSummaryTemplate, postID)
	if err := pc.Del(ctx, key).Err(); err != nil {
		return err
	}
	if err := pc.ZRem(ctx, GenerateRedisKey(PostTimeTemplate), postID).Err(); err != nil {
		return err
	}
	if err := pc.ZRem(ctx, GenerateRedisKey(PostRankingTemplate), postID).Err(); err != nil {
		return err
	}
	return nil
}

func (pc *PostCache) DeleteRedisPostSummary(ctx context.Context, postID int64) error {
	key := GenerateRedisKey(PostSummaryTemplate, postID)
	if err := pc.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func NewPostCache(client *redis.Client) PostCacheInterface {
	return &PostCache{
		client,
	}
}
