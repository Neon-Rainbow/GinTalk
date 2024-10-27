package kafka

import (
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/model"
	"context"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

// 使用 sync.Map 作为投票缓存
var voteMap sync.Map

// 聚合投票请求
func aggregateVote(postID int64, voteChange int) {
	// 更新投票计数
	value, _ := voteMap.LoadOrStore(postID, 0)
	voteCount := value.(int) + voteChange
	voteMap.Store(postID, voteCount)
}

// 定时批量更新数据库和 Redis
func (k *Kafka) flushAndUpdateVotes(ctx context.Context) {
	voteMap.Range(func(key, value any) bool {
		postID := key.(int64)
		voteChange := value.(int)

		// 调用更新函数
		err := updatePostVoteCount(ctx, k.voteDao, k.voteCache, postID, voteChange)
		if err != nil {
			zap.L().Error("更新投票失败", zap.Int64("postID", postID), zap.Error(err))
		} else {
			zap.L().Info("成功更新投票", zap.Int64("postID", postID), zap.Int("voteChange", voteChange))
		}

		// 删除已处理的记录
		voteMap.Delete(key)
		return true
	})
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
	err = updatePostVoteCount(ctx, k.voteDao, k.voteCache, id, 1)
	return nil
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
	err = updatePostVoteCount(ctx, k.voteDao, k.voteCache, id, -1)
	return nil
}

// updatePostVoteCount 更新帖子的投票数,同时更新 Redis 的帖子热度
func updatePostVoteCount(
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
