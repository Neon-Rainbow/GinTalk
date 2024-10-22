package service

import (
	"GinTalk/DTO"
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"GinTalk/pkg/snowflake"
	"context"
	"fmt"
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

	summary := TruncateByWords(postDTO.Content, MaxSummaryLength)

	err = ps.PostDaoInterface.CreatePost(ctx, postDTO, summary)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("保存帖子失败: %v", err),
		}
	}

	postSummary := &DTO.PostSummary{
		PostID:      postID,
		CommunityID: postDTO.CommunityID,
		Title:       postDTO.Title,
		AuthorId:    postDTO.AuthorId,
		Summary:     summary,
	}

	// 将帖子 ID 存入 Redis
	err = ps.PostCacheInterface.SavePostToRedis(ctx, postSummary)

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
	if postDTO.PostID == 0 {
		return &apiError.ApiError{
			Code: code.InvalidParam,
			Msg:  "postID 不能为空",
		}
	}

	fmt.Printf("截断前: %s\n", postDTO.Content)
	summary := TruncateByWords(postDTO.Content, MaxSummaryLength)
	fmt.Printf("截断后: %s\n", summary)

	err := ps.PostDaoInterface.UpdatePost(ctx, postDTO, summary)
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
