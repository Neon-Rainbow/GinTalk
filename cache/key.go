package cache

import "fmt"

const (
	BlackListTokenKeyTemplate = "blacklist:token:%v"

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
