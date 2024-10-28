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

// GenerateRedisKey 通过格式化给定的模板字符串和提供的参数生成一个 Redis key。
//
// 参数:
//   - template: 一个包含格式占位符的字符串模板。
//   - param: 一个可变参数，表示要格式化到模板中的值。
//
// 返回值:
//   - string: 一个格式化的字符串，可以用作 Redis key。
func GenerateRedisKey(template string, param ...any) string {
	return fmt.Sprintf(template, param...)
}
