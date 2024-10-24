package kafka

import (
	"GinTalk/DTO"
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/model"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

var kafkaHandleMap = map[string]func(ctx context.Context, k *Kafka, message KafkaMessage) error{
	MessageTypeAddPostVote:     addPostVoteHandle,
	MessageTypeSubPostVote:     subPostVoteHandle,
	MessageTypeSavePostToRedis: savePostToRedis,
}

// addPostVoteHandle 帖子投票数增加处理
func addPostVoteHandle(ctx context.Context, k *Kafka, message KafkaMessage) error {
	vote, err := handleKafkaMessage[model.KafkaVotePostModel](ctx, message)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(vote.PostId, 10, 64)
	if err != nil {
		zap.L().Error("Failed to convert PostId", zap.Error(err))
		return err
	}
	return UpdatePostVoteCount(ctx, k.voteDao, k.voteCache, id, 1)
}

// subPostVoteHandle 减少帖子投票数
func subPostVoteHandle(ctx context.Context, k *Kafka, message KafkaMessage) error {
	vote, err := handleKafkaMessage[model.KafkaVotePostModel](ctx, message)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(vote.PostId, 10, 64)
	if err != nil {
		zap.L().Error("Failed to convert PostId", zap.Error(err))
		return err
	}
	return UpdatePostVoteCount(ctx, k.voteDao, k.voteCache, id, -1)
}

// savePostToRedis 保存帖子到 Redis
func savePostToRedis(ctx context.Context, k *Kafka, message KafkaMessage) error {
	post, err := handleKafkaMessage[DTO.PostSummary](ctx, message)
	if err != nil {
		return err
	}
	if err := k.postCache.SavePostToRedis(ctx, &post); err != nil {
		zap.L().Error("Failed to save post to Redis", zap.Error(err))
		return err
	}
	zap.L().Info("Post saved to Redis successfully", zap.Any("post", post))
	return nil
}

// UpdatePostVoteCount 更新帖子的投票数,同时更新 Redis 的帖子热度
func UpdatePostVoteCount(
	ctx context.Context,
	voteDao dao.PostVoteDaoInterface,
	voteCache cache.VoteCacheInterface,
	postID int64,
	voteChange int,
) error {
	// 更新帖子的热度
	newDetail, err := voteDao.GetPostVoteCount(ctx, postID)
	if err != nil {
		zap.L().Error("Failed to get post vote count", zap.Error(err))
		return err
	}
	newUp := int(newDetail.Vote)
	oldUp := newUp - voteChange

	// err = v.VoteCacheInterface.UpdatePostHot(ctx, postID, newUp, createTime)
	err = voteCache.AddPostHot(ctx, postID, oldUp, newUp)
	if err != nil {
		zap.L().Error("Failed to update post hot score", zap.Error(err))
		return err
	}

	zap.L().Info("Update post vote count successfully", zap.Int64("postID", postID), zap.Int("vote", voteChange))
	return nil
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
