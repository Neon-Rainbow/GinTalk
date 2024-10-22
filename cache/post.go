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
	return sign*order + seconds/45000
}

func deltaHot(oldUp, newUp int) float64 {
	return math.Log10(max(float64(newUp), 1)) - math.Log10(max(float64(oldUp), 1))
}

type PostCacheInterface interface {
	SavePostToRedis(ctx context.Context, summary *DTO.PostSummary) error
	GetPostIDsFromRedis(ctx context.Context, order, pageNum, pageSize int) ([]int64, error)
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

func NewPostCache(client *redis.Client) PostCacheInterface {
	return &PostCache{
		client,
	}
}
