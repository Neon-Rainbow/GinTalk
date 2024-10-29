package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/kafka"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var voteGroup singleflight.Group
var postVoteCountGroup singleflight.Group

// VotePost 投票
// 投票使用kafka异步处理
// VotePost 处理用户对帖子的投票过程。
// 它使用 singleflight 机制确保投票操作仅执行一次，以防止重复投票。
//
// 参数:
//   - ctx: 请求的上下文，用于取消和截止日期。
//   - postID: 被投票的帖子的ID。
//   - userID: 投票用户的ID。
//
// 返回值:
//   - *apiError.ApiError: 如果投票过程失败，返回包含错误代码和消息的错误对象；
//     如果投票成功，则返回nil。
func VotePost(ctx context.Context, postID int64, userID int64) *apiError.ApiError {
	key := GenerateSingleFlightKey(SingleFlightKeyVotePost, postID, userID)
	go func() {
		_, err, _ := voteGroup.Do(key, func() (interface{}, error) {
			err := kafka.SendLikeMessage(ctx, &kafka.Vote{
				PostID: strconv.FormatInt(postID, 10),
				UserID: strconv.FormatInt(userID, 10),
				Vote:   1,
			})
			if err != nil {
				zap.L().Error("消息发送失败", zap.Error(err))
				return nil, err
			}
			zap.L().Info("投票成功")
			return nil, nil
		})
		if err != nil {
			zap.L().Error("投票操作失败", zap.Error(err))
		}
	}()
	return nil
}

// RevokeVotePost 处理用户对帖子的取消投票过程。
// 它使用 singleflight 机制确保取消投票操作仅执行一次，以防止重复操作。
//
// 参数:
//   - ctx: 请求的上下文，用于取消和截止日期。
//   - postID: 被取消投票的帖子的ID。
//   - userID: 取消投票用户的ID。
//
// 返回值:
//   - *apiError.ApiError: 如果取消投票过程失败，返回包含错误代码和消息的错误对象；
//     如果取消投票成功，则返回nil。
func RevokeVotePost(ctx context.Context, postID int64, userID int64) *apiError.ApiError {
	key := GenerateSingleFlightKey(SingleFlightKeyVotePost, postID, userID)
	go func() {
		_, err, _ := voteGroup.Do(key, func() (interface{}, error) {
			err := kafka.SendLikeMessage(ctx, &kafka.Vote{
				PostID: strconv.FormatInt(postID, 10),
				UserID: strconv.FormatInt(userID, 10),
				Vote:   0,
			})
			if err != nil {
				zap.L().Error("消息发送失败", zap.Error(err))
				return nil, err
			}
			zap.L().Info("撤销投票成功")
			return nil, nil
		})
		if err != nil {
			zap.L().Error("撤销投票操作失败", zap.Error(err))
		}
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
	key := GenerateSingleFlightKey(SingleFlightKeyPostVoteCount, postID)
	up, err, _ := postVoteCountGroup.Do(key, func() (interface{}, error) {
		up, err := dao.GetPostVoteCount(ctx, postID)
		if err != nil {
			return nil, &apiError.ApiError{
				Code: code.ServerError,
				Msg:  "查询错误",
			}
		}
		return up, nil
	})
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询错误",
		}
	}
	return up.(*DTO.PostVoteCounts), nil
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
