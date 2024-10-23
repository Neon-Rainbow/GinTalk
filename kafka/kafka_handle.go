package kafka

import (
	"GinTalk/model"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

var kafkaHandleMap = map[string]func(ctx context.Context, k *Kafka, message KafkaMessage) error{
	MessageTypeAddPostVote: addPostVoteHandle,
	MessageTypeSubPostVote: subPostVoteHandle,
}

func addPostVoteHandle(ctx context.Context, k *Kafka, message KafkaMessage) error {
	var vote model.KafkaVotePostModel
	msgMap, ok := message.Message.(map[string]interface{})
	if ok {
		msgBytes, err := json.Marshal(msgMap)
		if err != nil {
			zap.L().Error("Failed to marshal message to JSON", zap.Error(err))
			return err
		}
		if err := json.Unmarshal(msgBytes, &vote); err != nil {
			zap.L().Error("Failed to unmarshal JSON to VotePostDTO", zap.Error(err))
			return err
		}
	} else if msgBytes, ok := message.Message.([]byte); ok {
		// 如果消息是 []byte，直接反序列化
		if err := json.Unmarshal(msgBytes, &vote); err != nil {
			zap.L().Error("Failed to unmarshal byte message to VotePostDTO", zap.Error(err))
			return err
		}
	} else {
		zap.L().Error("Unsupported message type", zap.Any("message", message.Message))
		return fmt.Errorf("unsupported message type: %T", message.Message)
	}
	id, err := strconv.ParseInt(vote.PostId, 10, 64)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return err
	}
	err = k.voteDao.AddContentVoteUp(ctx, id)
	if err != nil {
		zap.L().Error("addPostVoteHandle.AddContentVoteUp failed", zap.Error(err))
		return err
	}
	zap.L().Info("addPostVoteHandle.AddContentVoteUp success", zap.Any("vote", vote))
	return nil
}

func subPostVoteHandle(ctx context.Context, k *Kafka, message KafkaMessage) error {
	var vote model.KafkaVotePostModel
	msgMap, ok := message.Message.(map[string]interface{})
	if ok {
		msgBytes, err := json.Marshal(msgMap)
		if err != nil {
			zap.L().Error("Failed to marshal message to JSON", zap.Error(err))
			return err
		}
		if err := json.Unmarshal(msgBytes, &vote); err != nil {
			zap.L().Error("Failed to unmarshal JSON to VotePostDTO", zap.Error(err))
			return err
		}
	} else if msgBytes, ok := message.Message.([]byte); ok {
		// 如果消息是 []byte，直接反序列化
		if err := json.Unmarshal(msgBytes, &vote); err != nil {
			zap.L().Error("Failed to unmarshal byte message to VotePostDTO", zap.Error(err))
			return err
		}
	} else {
		zap.L().Error("Unsupported message type", zap.Any("message", message.Message))
		return fmt.Errorf("unsupported message type: %T", message.Message)
	}
	id, err := strconv.ParseInt(vote.PostId, 10, 64)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return err
	}

	err = k.voteDao.SubContentVoteUp(ctx, id)
	if err != nil {
		zap.L().Error("subPostVoteHandle.SubContentVoteUp failed", zap.Error(err))
		return err
	}
	zap.L().Info("subPostVoteHandle.SubContentVoteUp success", zap.Any("vote", vote))
	return nil
}
