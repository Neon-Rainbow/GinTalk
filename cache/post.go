package cache

import (
	"GinTalk/DTO"
	"GinTalk/model"
	"context"
	"github.com/go-redis/redis/v8"
	"math"
	"time"
)

var _ PostCacheInterface = (*PostCache)(nil)

const (
	OrderByHot = iota
	OrderByTime
)

// hot 用于计算帖子的热度
// Reddit 热度算法
func hot(ups, downs int, date time.Time) float64 {
	s := float64(ups - downs)
	order := math.Log10(math.Max(math.Abs(s), 1))
	var sign float64
	if s > 0 {
		sign = 1
	} else if s == 0 {
		sign = 0
	} else {
		sign = -1
	}
	seconds := float64(date.Second() - 1577808000)
	return math.Round(sign*order + seconds/43200)
}

type PostCacheInterface interface {
	SavePostToRedis(ctx context.Context, summary *DTO.PostSummary) error
	GetPostInfoFromRedis(ctx context.Context, order int, pageNum, pageSize int) []DTO.PostSummary
	UpdatePostInfo(ctx context.Context, postID int64, post *model.Post) error
}

type PostCache struct {
	*redis.Client
}

func (pc *PostCache) SavePostToRedis(ctx context.Context, summary *DTO.PostSummary) error {
	//TODO implement me
	panic("implement me")
}

func (pc *PostCache) GetPostInfoFromRedis(ctx context.Context, order int, pageNum, pageSize int) []DTO.PostSummary {
	//TODO implement me
	panic("implement me")
}

func (pc *PostCache) UpdatePostInfo(ctx context.Context, postID int64, post *model.Post) error {
	//TODO implement me
	panic("implement me")
}

func NewPostCache(client *redis.Client) PostCacheInterface {
	return &PostCache{
		client,
	}
}
