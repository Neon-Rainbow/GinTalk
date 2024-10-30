// Package kafka 提供与 Kafka 消息系统交互的功能。
// 该包包括用于从 Kafka 主题中生产和消费消息的方法。
package kafka

const (
	// TopicCreatePost 帖子主题
	TopicCreatePost = "post"
	// TopicLike 点赞主题
	TopicLike = "like"
	// TopicComment 评论主题
	TopicComment = "comment"
	// TopicNotification 通知主题
	TopicNotification = "notification"
)
