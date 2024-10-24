package service

import (
	"GinTalk/DTO"
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/kafka"
	"GinTalk/model"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var _ VotePostServiceInterface = (*VoteService)(nil)

type VotePostServiceInterface interface {
	// VotePost 用于投票
	VotePost(ctx context.Context, postID int64, userID int64) *apiError.ApiError

	// RevokeVotePost 用于取消投票
	RevokeVotePost(ctx context.Context, postID int64, userID int64) *apiError.ApiError

	// MyVotePostList 用于查询用户投票过的帖子
	MyVotePostList(ctx context.Context, userID int64, pageNum, pageSize int) ([]int64, *apiError.ApiError)

	// GetVotePostCount 用于查询帖子的投票数量
	GetVotePostCount(ctx context.Context, postID int64) (*DTO.PostVoteCounts, *apiError.ApiError)

	// GetBatchPostVoteCount 该函数用于批量查询帖子的投票数量
	GetBatchPostVoteCount(ctx context.Context, postIDs []int64) ([]DTO.PostVoteCounts, *apiError.ApiError)

	// CheckUserPostVoted 批量查询用户是否投票过
	CheckUserPostVoted(ctx context.Context, postIDs []int64, userID int64) ([]DTO.UserVotePostRelationsDTO, *apiError.ApiError)

	// GetPostVoteDetail 获取帖子的投票详情
	GetPostVoteDetail(ctx context.Context, postID int64, pageNum, pageSize int) ([]DTO.UserVotePostDetailDTO, *apiError.ApiError)
}

type VoteService struct {
	dao.PostVoteDaoInterface
	cache.VoteCacheInterface
	kafka.KafkaInterface
}

func (v *VoteService) VotePost(ctx context.Context, postID int64, userID int64) *apiError.ApiError {
	isVoted, err := v.PostVoteDaoInterface.CheckUserVotedPost(ctx, postID, userID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票记录失败",
		}
	}
	if isVoted {
		return &apiError.ApiError{
			Code: code.InvalidParam,
			Msg:  "已经投过票",
		}
	}

	err = v.PostVoteDaoInterface.AddPostVote(ctx, postID, userID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "投票失败",
		}
	}

	msg := model.KafkaVotePostModel{
		PostId: fmt.Sprintf("%d", postID),
	}
	go func() {
		err := v.KafkaInterface.ProduceMessage(context.Background(), kafka.MessageTypeAddPostVote, msg)
		if err != nil {
			zap.L().Error("Failed to produce message", zap.Error(err))
		}
		zap.L().Info("投票成功")
	}()

	return nil
}

// RevokeVotePost 取消投票
func (v *VoteService) RevokeVotePost(ctx context.Context, postID int64, userID int64) *apiError.ApiError {
	isVoted, err := v.PostVoteDaoInterface.CheckUserVotedPost(ctx, postID, userID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票记录失败",
		}
	}
	if !isVoted {
		return &apiError.ApiError{
			Code: code.InvalidParam,
			Msg:  "未投过票",
		}
	}

	err = v.PostVoteDaoInterface.RevokePostVote(ctx, postID, userID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "取消投票失败",
		}
	}

	msg := model.KafkaVotePostModel{
		PostId: fmt.Sprintf("%d", postID),
	}
	go func() {
		err := v.KafkaInterface.ProduceMessage(context.Background(), kafka.MessageTypeSubPostVote, msg)
		if err != nil {
			zap.L().Error("Failed to produce message", zap.Error(err))
		}
	}()

	return nil
}

func (v *VoteService) MyVotePostList(ctx context.Context, userID int64, pageNum, pageSize int) ([]int64, *apiError.ApiError) {
	voteRecord, err := v.PostVoteDaoInterface.GetUserVoteList(ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票记录失败",
		}
	}
	return voteRecord, nil
}

func (v *VoteService) GetVotePostCount(ctx context.Context, postID int64) (*DTO.PostVoteCounts, *apiError.ApiError) {
	up, err := v.PostVoteDaoInterface.GetPostVoteCount(ctx, postID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询错误",
		}
	}
	return up, nil
}

func (v *VoteService) GetBatchPostVoteCount(ctx context.Context, postIDs []int64) ([]DTO.PostVoteCounts, *apiError.ApiError) {
	resp, err := v.PostVoteDaoInterface.GetBatchPostVoteCount(ctx, postIDs)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询错误",
		}
	}
	return resp, nil
}

// CheckUserPostVoted 批量查询用户是否投票过
func (v *VoteService) CheckUserPostVoted(ctx context.Context, postIDs []int64, userID int64) ([]DTO.UserVotePostRelationsDTO, *apiError.ApiError) {
	votes, err := v.PostVoteDaoInterface.CheckUserVoted(ctx, postIDs, userID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("批量查询投票记录失败: %v", err),
		}
	}
	return votes, nil
}

func (v *VoteService) GetPostVoteDetail(ctx context.Context, postID int64, pageNum, pageSize int) ([]DTO.UserVotePostDetailDTO, *apiError.ApiError) {
	voteDetails, err := v.PostVoteDaoInterface.GetPostVoteDetail(ctx, postID, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("查询投票详情失败: %v", err),
		}
	}
	return voteDetails, nil
}

func NewVoteService(voteDaoInterface dao.PostVoteDaoInterface, voteCacheInterface cache.VoteCacheInterface, kafkaInterface kafka.KafkaInterface) VotePostServiceInterface {
	return &VoteService{
		voteDaoInterface,
		voteCacheInterface,
		kafkaInterface,
	}
}

const (
	MaxRetries   = 3               // 最大重试次数
	InitialDelay = 2 * time.Second // 初始重试间隔
)
