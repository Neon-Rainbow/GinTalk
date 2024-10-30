package service

import (
	"GinTalk/dao"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"go.uber.org/zap"
)

func VoteComment(userID, commentID int64) *apiError.ApiError {
	err := dao.VoteComment(userID, commentID)
	// 如果是由于索引冲突的插入失败,则表示已经投过票
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "投票失败",
		}
	}
	go updateCommentVoteCount(commentID, 1)
	return nil
}

func RemoveVoteComment(userID, commentID int64) *apiError.ApiError {
	err := dao.RemoveVoteComment(userID, commentID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "取消投票失败",
		}
	}
	go updateCommentVoteCount(commentID, -1)
	return nil
}

func GetVoteComment(userID, commentID int64) (int, *apiError.ApiError) {
	count, err := dao.GetVoteComment(userID, commentID)
	if err != nil {
		return 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取投票数失败",
		}
	}
	return count, nil
}

func GetVoteCommentList(userID int64, commentIDs []int64) (map[int64]int, *apiError.ApiError) {
	result, err := dao.GetCommentVoteStatusList(userID, commentIDs)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取投票状态失败",
		}
	}
	return result, nil
}

func updateCommentVoteCount(commentID int64, count int) {
	var err error
	if count > 0 {
		err = dao.IncrCommentVoteCount(commentID)
	} else {
		err = dao.DecrCommentVoteCount(commentID)
	}
	if err != nil {
		zap.L().Error("更新评论投票数失败", zap.Int64("commentID", commentID), zap.Error(err))
		return
	}
}
