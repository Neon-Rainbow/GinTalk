package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"GinTalk/pkg/snowflake"
	"context"
)

// GetTopComments 获取帖子的顶级评论
func GetTopComments(ctx context.Context, postID int64, pageSize, pageNum int) ([]DTO.Comment, *apiError.ApiError) {
	comments, err := dao.GetTopComments(ctx, postID, pageSize, pageNum)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论失败",
		}
	}
	resp := make([]DTO.Comment, len(comments))
	for i, comment := range comments {
		resp[i] = DTO.Comment{
			CommentID:  comment.CommentID,
			PostID:     comment.PostID,
			AuthorID:   comment.AuthorID,
			AuthorName: comment.AuthorName,
			Content:    comment.Content,
		}
	}

	return resp, nil
}

// GetSubComments 获取帖子的子评论
func GetSubComments(ctx context.Context, postID, parentID int64, pageSize, pageNum int) ([]DTO.Comment, *apiError.ApiError) {
	comments, err := dao.GetSubComments(ctx, postID, parentID, pageSize, pageNum)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论失败",
		}
	}
	resp := make([]DTO.Comment, len(comments))
	for i, comment := range comments {
		resp[i] = DTO.Comment{
			CommentID:  comment.CommentID,
			PostID:     comment.PostID,
			AuthorID:   comment.AuthorID,
			AuthorName: comment.AuthorName,
			Content:    comment.Content,
		}
	}
	return resp, nil
}

// GetCommentByID 获取评论
func GetCommentByID(ctx context.Context, commentID int64) (*DTO.Comment, *apiError.ApiError) {
	comment, err := dao.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论失败",
		}
	}
	resp := &DTO.Comment{
		CommentID:  comment.CommentID,
		PostID:     comment.PostID,
		AuthorID:   comment.AuthorID,
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
	}
	return resp, nil
}

// CreateComment 创建评论
func CreateComment(ctx context.Context, comment *DTO.CreateCommentRequest) *apiError.ApiError {
	id, _ := snowflake.GetID()
	commentModel := &model.Comment{
		CommentID:  id,
		PostID:     comment.PostID,
		AuthorID:   comment.AuthorID,
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		Status:     1,
	}
	err := dao.CreateComment(ctx, commentModel, comment.ReplyID, comment.ParentID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "创建评论失败",
		}
	}
	return nil
}

// UpdateComment 更新评论
func UpdateComment(ctx context.Context, comment *DTO.Comment) *apiError.ApiError {
	err := dao.UpdateComment(ctx, comment.CommentID, comment.Content)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "更新评论失败",
		}
	}
	return nil
}

// DeleteComment 删除评论
func DeleteComment(ctx context.Context, commentID int64) *apiError.ApiError {
	err := dao.DeleteComment(ctx, commentID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "删除评论失败",
		}
	}
	return nil
}

// GetCommentCount 获取评论数量
func GetCommentCount(ctx context.Context, postID int64) (int64, *apiError.ApiError) {
	count, err := dao.GetCommentCount(ctx, postID)
	if err != nil {
		return 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论数量失败",
		}
	}
	return count, nil
}

// GetTopCommentCount 获取顶级评论数量
func GetTopCommentCount(ctx context.Context, postID int64) (int64, *apiError.ApiError) {
	count, err := dao.GetTopCommentCount(ctx, postID)
	if err != nil {
		return 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论数量失败",
		}
	}
	return count, nil
}

// GetSubCommentCount 获取子评论数量
func GetSubCommentCount(ctx context.Context, parentID int64) (int64, *apiError.ApiError) {
	count, err := dao.GetSubCommentCount(ctx, parentID)
	if err != nil {
		return 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论数量失败",
		}
	}
	return count, nil
}

func GetCommentCountByUserID(ctx context.Context, userID int64) (int64, *apiError.ApiError) {
	count, err := dao.GetCommentCountByUserID(ctx, userID)
	if err != nil {
		return 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论数量失败",
		}
	}
	return count, nil
}
