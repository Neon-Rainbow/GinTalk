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
func hot(ups, downs int, date time.Time) float64 {
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
	hotScore := hot(0, 0, time.Now())

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

// 目前不将帖子的内容存储到 Redis 中，只将 Redis 用于存储帖子的热度和时间

//func (pc *PostCache) GetPostInfoFromRedis(ctx context.Context, order, pageNum, pageSize int) ([]DTO.PostSummary, error) {
//	var key string
//	if order == OrderByHot {
//		key = GenerateRedisKey(PostRankingTemplate)
//	} else {
//		key = GenerateRedisKey(PostTimeTemplate)
//	}
//
//	// 计算分页的开始和结束位置
//	start := int64((pageNum - 1) * pageSize)
//	end := start + int64(pageSize) - 1
//
//	// 从 Redis 有序集合中获取帖子 ID 列表
//	postIDs, err := pc.ZRevRange(ctx, key, start, end).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	// 根据获取的 ID 列表逐个获取帖子内容
//	var posts []DTO.PostSummary
//	for _, id := range postIDs {
//		data, err := pc.Get(ctx, GenerateRedisKey(PostSummaryTemplate, id)).Result()
//		if err != nil {
//			continue // 忽略无法获取的帖子
//		}
//		var post DTO.PostSummary
//		if err := json.Unmarshal([]byte(data), &post); err != nil {
//			continue
//		}
//		posts = append(posts, post)
//	}
//	return posts, nil
//}
//
//func (pc *PostCache) UpdatePostInfo(ctx context.Context, postID int64, post *DTO.PostSummary) error {
//	key := GenerateRedisKey(PostSummaryTemplate, post.PostID)
//	data, err := json.Marshal(post)
//	if err != nil {
//		return err
//	}
//	// 更新帖子数据
//	err = pc.Set(ctx, key, data, PostStoreTime).Err()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//

func NewPostCache(client *redis.Client) PostCacheInterface {
	return &PostCache{
		client,
	}
}
