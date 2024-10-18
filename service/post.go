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
	CreatePost(ctx context.Context, postDTO *DTO.PostDetail) *apiError.ApiError
	GetPostList(ctx context.Context, pageNum int, pageSize int) ([]*DTO.PostDetail, *apiError.ApiError)
	GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDetail, *apiError.ApiError)
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
	err = ps.postDao.WithContext(ctx).Create(&post)
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
	var list []*DTO.PostDetail
	/*
		SELECT `post`.`post_id`,`post`.`title`,`post`.`content`,`post`.`author_id`,`user`.`username`,`post`.`community_id`,`community`.`community_name`,`post`.`status`
		FROM `post`
		    INNER JOIN `community`
		        ON `community`.`community_id` = `post`.`community_id`
		    INNER JOIN `user`
		        ON `user`.`user_id` = `post`.`author_id`
		WHERE `post`.`delete_time` IS NULL LIMIT 10
	*/
	err := ps.postDao.WithContext(ctx).
		Limit(pageSize).Offset((pageNum-1)*pageSize).
		Select(dao.Post.PostID, dao.Post.Title, dao.Post.Content, dao.Post.AuthorID, dao.User.Username, dao.Post.CommunityID, dao.Community.CommunityName, dao.Post.Status).
		Join(dao.Community, dao.Community.CommunityID.EqCol(dao.Post.CommunityID)).
		Join(dao.User, dao.User.UserID.EqCol(dao.Post.AuthorID)).
		Scan(&list)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取帖子列表失败: %v", err),
		}
	}
	return list, nil
}

func (ps *PostService) GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDetail, *apiError.ApiError) {
	var postDetail DTO.PostDetail
	// 查询帖子详情
	/*
		SELECT
		    `post`.`post_id`,
		    `post`.`title`,
		    `post`.`content`,
		    `post`.`author_id`,
		    `user`.`username`,  -- 确保这个字段在 user 表中存在
		    `post`.`community_id`,
		    `community`.`community_name`,
		    `post`.`status`
		FROM
		    `post`
		INNER JOIN
		    `community` ON `post`.`community_id` = `community`.`community_id`
		INNER JOIN
		    `user` ON `post`.`author_id` = `user`.`user_id`
		WHERE
		    `post`.`post_id` = 254027605567602689
		    AND `post`.`delete_time` IS NULL
	*/
	err := ps.postDao.WithContext(ctx).
		Where(dao.Post.PostID.Eq(postID)).
		Select(dao.Post.PostID, dao.Post.Title, dao.Post.Content, dao.Post.AuthorID, dao.User.Username, dao.Post.CommunityID, dao.Community.CommunityName, dao.Post.Status).
		Join(dao.Community, dao.Community.CommunityID.EqCol(dao.Post.CommunityID)).
		Join(dao.User, dao.User.UserID.EqCol(dao.Post.AuthorID)).
		Scan(&postDetail)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取帖子详情失败: %v", err),
		}
	}
	return &postDetail, nil
}
