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

type PostServiceInterface interface {
	CreatePost(ctx context.Context, postDTO *DTO.PostDTO) *apiError.ApiError
	GetPostList(ctx context.Context, pageNum int, pageSize int) ([]*DTO.PostDTO, *apiError.ApiError)
	GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDTO, *apiError.ApiError)
}

type PostService struct {
	postDao dao.IPostDo
}

// NewPostService 使用依赖注入初始化 PostService
func NewPostService(postDao dao.IPostDo) PostServiceInterface {
	return &PostService{
		postDao: postDao,
	}
}

func (ps *PostService) CreatePost(ctx context.Context, postDTO *DTO.PostDTO) *apiError.ApiError {
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
	err = ps.postDao.WithContext(ctx).Create(&post)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("保存帖子失败: %v", err),
		}
	}
	return nil
}

func (ps *PostService) GetPostList(ctx context.Context, pageNum int, pageSize int) ([]*DTO.PostDTO, *apiError.ApiError) {
	list, err := ps.postDao.WithContext(ctx).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find()
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取帖子列表失败: %v", err),
		}
	}
	var postDTOList []*DTO.PostDTO
	err = copier.Copy(postDTOList, list)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("拷贝结构体失败: %v", err),
		}
	}
	return postDTOList, nil
}

func (ps *PostService) GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDTO, *apiError.ApiError) {
	postDetail, err := ps.postDao.WithContext(ctx).Where(dao.Post.PostID.Eq(postID)).First()
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取帖子详情失败: %v", err),
		}
	}
	var postDTO *DTO.PostDTO
	err = copier.Copy(postDTO, postDetail)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("拷贝结构体失败: %v", err),
		}
	}
	return postDTO, nil
}
