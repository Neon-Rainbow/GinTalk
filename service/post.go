package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"GinTalk/pkg/snowflake"
	"context"
	"fmt"
	"github.com/jinzhu/copier"
)

var _ PostServiceInterface = (*PostService)(nil)

type PostServiceInterface interface {
	CreatePost(ctx context.Context, postDTO *DTO.PostDetail) *apiError.ApiError
	GetPostList(ctx context.Context, pageNum int, pageSize int) ([]*DTO.PostDetail, *apiError.ApiError)
	GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDetail, *apiError.ApiError)
}

type PostService struct {
	dao.IPostDo
	dao.PostDaoInterface
}

// NewPostService 使用依赖注入初始化 PostService
func NewPostService(postDao dao.IPostDo, postDaoInterface dao.PostDaoInterface) PostServiceInterface {
	return &PostService{
		IPostDo:          postDao,
		PostDaoInterface: postDaoInterface,
	}
}

func (ps *PostService) CreatePost(ctx context.Context, postDTO *DTO.PostDetail) *apiError.ApiError {
	postID, err := snowflake.GetID()
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("生成帖子ID失败: %v", err),
		}
	}
	postDTO.PostID = postID
	var post model.Post
	err = copier.Copy(&post, postDTO)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("拷贝结构体失败: %v", err),
		}
	}
	err = ps.PostDaoInterface.CreatePost(ctx, &post)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("保存帖子失败: %v", err),
		}
	}
	return nil
}

func (ps *PostService) GetPostList(ctx context.Context, pageNum int, pageSize int) ([]*DTO.PostDetail, *apiError.ApiError) {
	// pageNum 和 pageSize 不能小于等于 0
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	list, err := ps.PostDaoInterface.GetPostList(ctx, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取帖子列表失败: %v", err),
		}
	}
	return list, nil
}

func (ps *PostService) GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDetail, *apiError.ApiError) {
	postDetail, err := ps.PostDaoInterface.GetPostDetail(ctx, postID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取帖子详情失败: %v", err),
		}
	}
	return postDetail, nil
}
