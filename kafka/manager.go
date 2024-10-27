package kafka

import (
	"GinTalk/DTO"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
)

var manager *Manager

// Manager 管理多个生产者和消费者
type Manager struct {
	Brokers []string
	Writers map[string]*kafka.Writer
	Readers map[string]*kafka.Reader
}

// newKafkaManager 初始化 Manager
func newKafkaManager(brokers []string, topics []string, groupID string) *Manager {
	writers := make(map[string]*kafka.Writer)
	readers := make(map[string]*kafka.Reader)

	// 初始化生产者
	for _, topic := range topics {
		writers[topic] = kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		})
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

// sendMessage 发送消息到指定的Topic
func (km *Manager) sendMessage(ctx context.Context, topic string, key, value []byte) error {
	writer, exists := km.Writers[topic]
	if !exists {
		log.Printf("Producer for topic %s not found", topic)
		return nil
	}
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}
	if err := writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("Failed to send message to %s: %v", topic, err)
		return err
	}
	return nil
}

func SendPostMessage(ctx context.Context, postMsg *DTO.PostDetail) error {
	topic := TopicPost
	value, err := json.Marshal(postMsg)
	if err != nil {
		return err
	}
	return GetKafkaManager().sendMessage(ctx, topic, nil, value)
}

func SendLikeMessage(ctx context.Context, vote *Vote) error {
	topic := TopicLike
	value, err := json.Marshal(vote)
	if err != nil {
		return err
	}
	return GetKafkaManager().sendMessage(ctx, topic, nil, value)
}

func SendCommentMessage(ctx context.Context, commentMsg *DTO.CommentDetail) error {
	topic := TopicComment
	value, err := json.Marshal(commentMsg)
	if err != nil {
		return err
	}
	return GetKafkaManager().sendMessage(ctx, topic, nil, value)
}

// startConsuming 启动消费者消费指定Topic的消息
func (km *Manager) startConsuming(ctx context.Context, topic string) {
	reader, exists := km.Readers[topic]
	if !exists {
		log.Printf("Consumer for topic %s not found", topic)
		return
	}
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message from %s: %v", topic, err)
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

// InitKafkaManager 初始化 KafkaManager
func InitKafkaManager() {

	brokers := []string{"localhost:9092"}
	topics := []string{TopicPost, TopicLike, TopicComment}

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
	TopicPost:    handlePostMessage,
	TopicLike:    handleLikeMessage,
	TopicComment: handleCommentMessage,
}
