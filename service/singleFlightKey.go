package service

import (
	"fmt"
)

const (
	// SingleFlightKeyPostList 用于获取帖子列表的单飞模式 key, 三个参数分别为 order, pageNum, pageSize
	SingleFlightKeyPostList = "post_list_%d_%d_%d"

	// SingleFlightKeyPostDetail 用于获取帖子详情的单飞模式 key, 一个参数为 postID
	SingleFlightKeyPostDetail = "post_detail_%d"

	// SingleFlightKeyVotePost 用于投票的单飞模式 key, 两个参数分别为 postID, userID
	SingleFlightKeyVotePost = "vote_post_%d_%d"

	// SingleFlightKeyPostVoteCount 用于获取帖子投票数的单飞模式 key, 一个参数为 postID
	SingleFlightKeyPostVoteCount = "post_vote_count_%d"
)

// GenerateSingleFlightKey 通过使用提供的模板字符串和给定的参数生成单飞操作的唯一键。
//
// 参数:
//   - template: 定义键格式的字符串模板。
//   - params: 要格式化到模板中的可变参数列表。
//
// 返回值:
//   - key: 一个格式化的字符串，作为单飞操作的唯一键。
func GenerateSingleFlightKey(template string, params ...interface{}) (key string) {
	return fmt.Sprintf(template, params...)
}
