package service

import (
	"GinTalk/DTO"
	"GinTalk/cache"
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
	GetPostList(ctx context.Context, pageNum int, pageSize int, order int) ([]DTO.PostSummary, *apiError.ApiError)
	GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDetail, *apiError.ApiError)
	UpdatePost(ctx context.Context, postDTO *DTO.PostDetail) *apiError.ApiError
	GetPostListByCommunityID(ctx context.Context, communityID int64, pageNum int, pageSize int) ([]DTO.PostSummary, *apiError.ApiError)
}

type PostService struct {
	dao.PostDaoInterface
	cache.PostCacheInterface
}

// NewPostService 使用依赖注入初始化 PostService
func NewPostService(postDaoInterface dao.PostDaoInterface, cache cache.PostCacheInterface) PostServiceInterface {
	return &PostService{
		PostDaoInterface:   postDaoInterface,
		PostCacheInterface: cache,
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

	summary := TruncateByWords(post.Content, 100)
	post.Summary = summary

	err = ps.PostDaoInterface.CreatePost(ctx, &post)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("保存帖子失败: %v", err),
		}
	}

	err = ps.PostCacheInterface.SavePostToRedis(ctx, &DTO.PostSummary{
		PostID:      post.PostID,
		Title:       post.Title,
		Summary:     summary,
		AuthorId:    post.AuthorID,
		CommunityID: post.CommunityID,
	})

	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("保存帖子失败: %v", err),
		}
	}
	return nil
}

func (ps *PostService) GetPostList(ctx context.Context, pageNum int, pageSize int, order int) ([]DTO.PostSummary, *apiError.ApiError) {
	// pageNum 和 pageSize 不能小于等于 0
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	postIDs, err := ps.PostCacheInterface.GetPostIDsFromRedis(ctx, order, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取帖子列表失败: %v", err),
		}
	}

	list, err := ps.PostDaoInterface.GetPostListBatch(ctx, postIDs)
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

func (ps *PostService) UpdatePost(ctx context.Context, postDTO *DTO.PostDetail) *apiError.ApiError {
	post := model.Post{
		PostID:  postDTO.PostID,
		Title:   postDTO.Title,
		Summary: TruncateByWords(postDTO.Content, 100),
		Content: postDTO.Content,
	}

	err := ps.PostDaoInterface.UpdatePost(ctx, &post)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("更新帖子失败: %v", err),
		}
	}

	return nil
}

func (ps *PostService) GetPostListByCommunityID(ctx context.Context, communityID int64, pageNum int, pageSize int) ([]DTO.PostSummary, *apiError.ApiError) {
	// pageNum 和 pageSize 不能小于等于 0
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	list, err := ps.PostDaoInterface.GetPostListByCommunityID(ctx, communityID, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取社区帖子列表失败: %v", err),
		}
	}
	return list, nil
}
