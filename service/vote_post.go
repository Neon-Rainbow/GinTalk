package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/kafka"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"context"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

func VotePost(ctx context.Context, postID int64, userID int64) *apiError.ApiError {
	isVoted, err := dao.CheckUserVotedPost(ctx, postID, userID)
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

	err = dao.AddPostVote(ctx, postID, userID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "投票失败",
		}
	}

	go func() {
		err := kafka.SendLikeMessage(context.Background(), &kafka.Vote{PostID: strconv.FormatInt(postID, 10), UserID: strconv.FormatInt(userID, 10), Vote: 1})
		if err != nil {
			zap.L().Error("Failed to produce message", zap.Error(err))
		}
		zap.L().Info("投票成功")
	}()

	return nil
}

// RevokeVotePost 取消投票
func RevokeVotePost(ctx context.Context, postID int64, userID int64) *apiError.ApiError {
	isVoted, err := dao.CheckUserVotedPost(ctx, postID, userID)
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

	err = dao.RevokePostVote(ctx, postID, userID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "取消投票失败",
		}
	}

	go func() {
		err := kafka.SendLikeMessage(context.Background(), &kafka.Vote{PostID: strconv.FormatInt(postID, 10), UserID: strconv.FormatInt(userID, 10), Vote: -1})
		if err != nil {
			zap.L().Error("Failed to produce message", zap.Error(err))
		}
		zap.L().Info("取消投票成功")
	}()

	return nil
}

func MyVotePostList(ctx context.Context, userID int64, pageNum, pageSize int) ([]int64, *apiError.ApiError) {
	voteRecord, err := dao.GetUserVoteList(ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票记录失败",
		}
	}
	return voteRecord, nil
}

func GetVotePostCount(ctx context.Context, postID int64) (*DTO.PostVoteCounts, *apiError.ApiError) {
	up, err := dao.GetPostVoteCount(ctx, postID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询错误",
		}
	}
	return up, nil
}

func GetBatchPostVoteCount(ctx context.Context, postIDs []int64) ([]DTO.PostVoteCounts, *apiError.ApiError) {
	resp, err := dao.GetBatchPostVoteCount(ctx, postIDs)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询错误",
		}
	}
	return resp, nil
}

// CheckUserPostVoted 批量查询用户是否投票过
func CheckUserPostVoted(ctx context.Context, postIDs []int64, userID int64) ([]DTO.UserVotePostRelationsDTO, *apiError.ApiError) {
	votes, err := dao.CheckUserVoted(ctx, postIDs, userID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("批量查询投票记录失败: %v", err),
		}
	}
	return votes, nil
}

func GetPostVoteDetail(ctx context.Context, postID int64, pageNum, pageSize int) ([]DTO.UserVotePostDetailDTO, *apiError.ApiError) {
	voteDetails, err := dao.GetPostVoteDetail(ctx, postID, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("查询投票详情失败: %v", err),
		}
	}
	return voteDetails, nil
}
