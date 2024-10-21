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

	// PostRanking 在redis中存储帖子的热度
	PostRanking = "post:ranking:id:%v"
)

// GenerateRedisKey 生成redis key
func GenerateRedisKey(template string, param ...any) string {
	return fmt.Sprintf(template, param...)
}
