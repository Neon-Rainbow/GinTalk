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

type CommentServiceInterface interface {
	GetTopComments(ctx context.Context, postID int64, pageSize, pageNum int) ([]DTO.Comment, *apiError.ApiError)
	GetSubComments(ctx context.Context, postID, parentID int64, pageSize, pageNum int) ([]DTO.Comment, *apiError.ApiError)
	GetCommentByID(ctx context.Context, commentID int64) (*DTO.Comment, *apiError.ApiError)
	CreateComment(ctx context.Context, comment *DTO.CreateCommentRequest) *apiError.ApiError
	UpdateComment(ctx context.Context, comment *DTO.Comment) *apiError.ApiError
	DeleteComment(ctx context.Context, commentID int64) *apiError.ApiError
	GetCommentCount(ctx context.Context, postID int64) (int64, *apiError.ApiError)
	GetTopCommentCount(ctx context.Context, postID int64) (int64, *apiError.ApiError)
	GetSubCommentCount(ctx context.Context, parentID int64) (int64, *apiError.ApiError)
	GetCommentCountByUserID(ctx context.Context, userID int64) (int64, *apiError.ApiError)
}

type CommentService struct {
	dao.CommentDaoInterface
}

func NewCommentService(dao dao.CommentDaoInterface) CommentServiceInterface {
	return &CommentService{dao}
}

func (cs *CommentService) GetTopComments(ctx context.Context, postID int64, pageSize, pageNum int) ([]DTO.Comment, *apiError.ApiError) {
	comments, err := cs.CommentDaoInterface.GetTopComments(ctx, postID, pageSize, pageNum)
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

func (cs *CommentService) GetSubComments(ctx context.Context, postID, parentID int64, pageSize, pageNum int) ([]DTO.Comment, *apiError.ApiError) {
	comments, err := cs.CommentDaoInterface.GetSubComments(ctx, postID, parentID, pageSize, pageNum)
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

func (cs *CommentService) GetCommentByID(ctx context.Context, commentID int64) (*DTO.Comment, *apiError.ApiError) {
	comment, err := cs.CommentDaoInterface.GetCommentByID(ctx, commentID)
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

func (cs *CommentService) CreateComment(ctx context.Context, comment *DTO.CreateCommentRequest) *apiError.ApiError {
	id, _ := snowflake.GetID()
	commentModel := &model.Comment{
		CommentID:  id,
		PostID:     comment.PostID,
		AuthorID:   comment.AuthorID,
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		Status:     1,
	}
	err := cs.CommentDaoInterface.CreateComment(ctx, commentModel, comment.ReplyID, comment.ParentID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "创建评论失败",
		}
	}
	return nil
}

func (cs *CommentService) UpdateComment(ctx context.Context, comment *DTO.Comment) *apiError.ApiError {
	err := cs.CommentDaoInterface.UpdateComment(ctx, comment.CommentID, comment.Content)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "更新评论失败",
		}
	}
	return nil
}

func (cs *CommentService) DeleteComment(ctx context.Context, commentID int64) *apiError.ApiError {
	err := cs.CommentDaoInterface.DeleteComment(ctx, commentID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "删除评论失败",
		}
	}
	return nil
}

func (cs *CommentService) GetCommentCount(ctx context.Context, postID int64) (int64, *apiError.ApiError) {
	count, err := cs.CommentDaoInterface.GetCommentCount(ctx, postID)
	if err != nil {
		return 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论数量失败",
		}
	}
	return count, nil
}

func (cs *CommentService) GetTopCommentCount(ctx context.Context, postID int64) (int64, *apiError.ApiError) {
	count, err := cs.CommentDaoInterface.GetTopCommentCount(ctx, postID)
	if err != nil {
		return 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论数量失败",
		}
	}
	return count, nil
}

func (cs *CommentService) GetSubCommentCount(ctx context.Context, parentID int64) (int64, *apiError.ApiError) {
	count, err := cs.CommentDaoInterface.GetSubCommentCount(ctx, parentID)
	if err != nil {
		return 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论数量失败",
		}
	}
	return count, nil
}

func (cs *CommentService) GetCommentCountByUserID(ctx context.Context, userID int64) (int64, *apiError.ApiError) {
	count, err := cs.CommentDaoInterface.GetCommentCountByUserID(ctx, userID)
	if err != nil {
		return 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取评论数量失败",
		}
	}
	return count, nil
}
