package kafka

import (
	"GinTalk/DTO"
	"context"
	"go.uber.org/zap"
)

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
