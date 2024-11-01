package websocket

const (
	// NotificationTypeComment 评论通知
	NotificationTypeComment = iota

	// NotificationTypeVote 点赞通知
	NotificationTypeVote
)

const (
	// MessageKindText 文本消息
	MessageKindText = "text"

	// MessageKindOnline 上线消息
	MessageKindOnline = "online"

	// MessageKindOffline 下线消息
	MessageKindOffline = "offline"

	// MessageKindNotificationComment 评论通知
	MessageKindNotificationComment = "notification_comment"

	// MessageKindNotificationVote 点赞通知
	MessageKindNotificationVote = "notification_vote"
)

// Message 是 websocket 传输的消息
type Message struct {
	// Kind 消息类型
	Kind string `json:"kind"`

	// From 发送者
	From string `json:"from"`

	// To 接收者
	To string `json:"to,omitempty"`

	// Data 消息内容
	Data string `json:"data,omitempty"`
}
