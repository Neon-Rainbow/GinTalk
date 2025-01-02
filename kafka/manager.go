package kafka

import (
	"GinTalk/DTO"
	"GinTalk/settings"
	"GinTalk/websocket"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var manager *Manager
var once sync.Once

// Manager 管理多个生产者和消费者
type Manager struct {
	Brokers []string
	Writers map[string]*kafka.Writer
	Readers map[string]*kafka.Reader
}

// newKafkaManager 使用提供的 brokers、topics 和 group ID 初始化一个新的 Kafka 管理器。
// 它为每个 topic 创建 Kafka writers 以处理消息生产，并为每个 topic 创建 Kafka readers 以处理消息消费。
//
// 参数:
//   - brokers([]string): 表示 Kafka broker 地址的字符串切片。
//   - topics([]string): 表示要管理的 Kafka 主题的字符串切片。
//   - groupID(string): 表示 Kafka 消费者组 ID 的字符串。
//
// 返回值:
//   - *Manager: 一个指向包含已初始化 Kafka writers 和 readers 的 Manager 结构体的指针。
func newKafkaManager(brokers []string, topics []string, groupID string) *Manager {
	writers := make(map[string]*kafka.Writer)
	readers := make(map[string]*kafka.Reader)

	// 初始化之前先验证 Kafka 是否可以连接
	for _, broker := range brokers {
		conn, err := kafka.Dial("tcp", broker)
		if err != nil {
			zap.L().Fatal("连接 Kafka 失败", zap.Error(err))
			return nil
		}
		conn.Close()
	}

	// 创建 Kafka 主题
	for _, topic := range topics {
		conn, err := kafka.Dial("tcp", brokers[0])
		if err != nil {
			zap.L().Fatal("连接 Kafka 失败", zap.Error(err))
			return nil
		}
		if err := conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}); err != nil {
			zap.L().Fatal("创建 Kafka 主题失败", zap.Error(err))
		}
		conn.Close()
	}

	// 初始化生产者
	for _, topic := range topics {
		writers[topic] = &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
	}

	// 初始化消费者
	for _, topic := range topics {
		readers[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})
	}

	return &Manager{Brokers: brokers, Writers: writers, Readers: readers}
}

// sendMessage 将消息发送到指定的 Kafka 主题。
// 它从 Manager 的 Writers 映射中检索给定主题的 Kafka writer。
// 如果 writer 不存在，它会记录一条消息并返回 nil。
// 如果 writer 存在，它会构造一个带有提供的 key 和 value 的 kafka.Message，
// 并尝试使用 writer 写入消息。如果写入操作失败，
// 它会记录错误并返回该错误。
//
// 参数:
//   - ctx: 写入操作的上下文。
//   - topic: 要发送消息的 Kafka 主题。
//   - key: 消息的键。
//   - value: 消息的值。
//
// 返回值:
//   - error: 如果消息无法发送，则返回错误，否则返回 nil。
func (km *Manager) sendMessage(ctx context.Context, topic string, key, value []byte) error {
	writer, exists := km.Writers[topic]
	if !exists {
		zap.L().Error("生产者不存在", zap.String("topic", topic))
		return nil
	}
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}
	if err := writer.WriteMessages(ctx, msg); err != nil {
		zap.L().Error("发送消息失败", zap.Error(err))
		return err
	}
	return nil
}

// SendPostMessage 向 Kafka 主题发送帖子消息。
// 它将提供的帖子消息序列化为 JSON 格式并发送到 Kafka 管理器。
//
// 参数:
//   - ctx: 控制取消和截止日期的上下文。
//   - postMsg: 指向包含帖子消息详细信息的 PostDetail DTO 的指针。
//
// 返回值:
//   - error: 如果消息无法序列化或发送，则返回错误。
func SendPostMessage(ctx context.Context, postMsg *DTO.PostDetail) error {
	topic := TopicCreatePost
	value, err := json.Marshal(postMsg)
	if err != nil {
		return err
	}
	return GetKafkaManager().sendMessage(ctx, topic, nil, value)
}

// SendLikeMessage 发送点赞消息到 Kafka 主题。
// 它将提供的 Vote 对象序列化为 JSON 格式并使用 Kafka 管理器发送。
//
// 参数:
//   - ctx: 控制取消和截止日期的上下文。
//   - vote: 指向要发送的 Vote 对象的指针。
//
// 返回值:
//   - error: 如果消息无法发送，则返回错误，否则返回 nil。
func SendLikeMessage(ctx context.Context, vote *Vote) error {
	topic := TopicLike
	value, err := json.Marshal(vote)
	if err != nil {
		return err
	}
	return GetKafkaManager().sendMessage(ctx, topic, nil, value)
}

// SendCommentMessage 发送评论消息到 Kafka 主题。
// 它将提供的评论消息序列化为 JSON 格式并使用 Kafka 管理器发送。
//
// 参数:
//   - ctx: 控制取消和截止日期的上下文。
//   - commentMsg: 指向包含评论消息详细信息的 CommentDetail 结构体的指针。
//
// 返回值:
//   - error: 如果消息无法序列化或发送，则返回错误。
func SendCommentMessage(ctx context.Context, commentMsg *DTO.CommentDetail) error {
	topic := TopicComment
	value, err := json.Marshal(commentMsg)
	if err != nil {
		return err
	}
	return GetKafkaManager().sendMessage(ctx, topic, nil, value)
}

// SendNotificationMessage 发送通知消息到 Kafka 主题。
// 它将提供的通知消息序列化为 JSON 格式并使用 Kafka 管理器发送。
//
// 参数:
//   - ctx: 控制取消和截止日期的上下文。
//   - notificationMsg: 指向要发送的通知消息的指针。
//
// 返回值:
//   - error: 如果消息无法发送，则返回错误，否则返回 nil。
func SendNotificationMessage(ctx context.Context, notificationMsg *websocket.Message) error {
	return fmt.Errorf("函数以及被废弃")
}

// startConsuming 启动消费者消费指定Topic的消息
func (km *Manager) startConsuming(ctx context.Context, topic string) {
	reader, exists := km.Readers[topic]
	if !exists {
		zap.L().Error("消费者不存在", zap.String("topic", topic))
		return
	}
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			zap.L().Error("读取消息失败", zap.Error(err))
			break
		}
		handles[topic](msg)
	}
}

// Close 关闭所有生产者和消费者
func (km *Manager) Close() {
	for _, writer := range km.Writers {
		err := writer.Close()
		if err != nil {
			zap.L().Error("关闭生产者失败", zap.Error(err))
		}
	}
	for _, reader := range km.Readers {
		err := reader.Close()
		if err != nil {
			zap.L().Error("关闭消费者失败", zap.Error(err))
		}
	}
}

// InitKafkaManager 初始化 KafkaManager 单例实例。
// 它设置 Kafka brokers 和 topics，并在单独的 goroutine 中开始消费指定 topics 的消息。
// 此函数使用 sync.Once 机制确保初始化只执行一次。
func InitKafkaManager() {
	brokers := settings.GetConfig().KafkaConfig.Brokers
	topics := []string{TopicCreatePost, TopicLike, TopicComment, TopicNotification}

	// 初始化 KafkaManager
	manager = newKafkaManager(brokers, topics, "example-group")

	for _, topic := range topics {
		go manager.startConsuming(context.Background(), topic)
	}
}

func GetKafkaManager() *Manager {
	return manager
}

type handleFunc func(kafka.Message)

var handles = map[string]handleFunc{
	TopicCreatePost: handleCreatePostMessage,
	TopicLike:       handleLikeMessage,
	TopicComment:    handleCommentMessage,
}
