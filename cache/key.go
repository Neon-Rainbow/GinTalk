package cache

import "fmt"

const (
	// UserTokenKeyTemplate 用户令牌 key 模板
	UserTokenKeyTemplate = "user:token:%v"

	BlackListTokenKeyTemplate = "blacklist:token:%v"

	// CommunityDetailTemplate 社区详情 key 模板
	CommunityDetailTemplate = "community:id:%v"

	// CommunityListTemplate 社区列表 key 模板
	CommunityListTemplate = "community:list"

	// PostSummaryTemplate 用于在 redis 中存储帖子的概述信息
	PostSummaryTemplate = "post:id:%v"

	// PostRankingTemplate 在redis中存储帖子的热度
	PostRankingTemplate = "post:ranking"

	// PostTimeTemplate 在 Redis 中存储帖子的时间
	PostTimeTemplate = "post:time"
)

// GenerateRedisKey 生成redis key
func GenerateRedisKey(template string, param ...any) string {
	return fmt.Sprintf(template, param...)
}
