package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

var kafkaHandleMap = map[string]func(ctx context.Context, k *Kafka, message KafkaMessage) error{
	MessageTypeAddPostVote:     addPostVoteHandle,
	MessageTypeSubPostVote:     subPostVoteHandle,
	MessageTypeSavePostToRedis: savePostToRedis,
}

// handleKafkaMessage 处理 Kafka 消息,将消息解析为泛型类型 T
func handleKafkaMessage[T any](ctx context.Context, message KafkaMessage) (T, error) {
	var result T
	var msgBytes []byte
	var err error

	// 将消息统一转换为字节数组
	switch msg := message.Message.(type) {
	case map[string]interface{}:
		msgBytes, err = json.Marshal(msg)
		if err != nil {
			zap.L().Error("无法将消息封送至 JSON", zap.Error(err))
			return result, err
		}
	case []byte:
		msgBytes = msg
	default:
		err = fmt.Errorf("不支持的消息类型: %T", message.Message)
		zap.L().Error("不支持的消息类型", zap.Any("message", message.Message))
		return result, err
	}

	// 解析字节数组为泛型类型 T
	if err := json.Unmarshal(msgBytes, &result); err != nil {
		zap.L().Error("无法解组消息", zap.Error(err))
		return result, err
	}

	return result, nil
}
