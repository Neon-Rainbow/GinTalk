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

func GenerateSingleFlightKey(template string, params ...interface{}) string {
	return fmt.Sprintf(template, params...)
}
