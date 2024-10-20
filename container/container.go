// Package container: 依赖注入容器，用于初始化 Service 实例。
package container

import (
	"GinTalk/dao"
	"GinTalk/dao/MySQL"
	"GinTalk/service"
	"context"
	"sync"
)

var (
	once             sync.Once
	communityService service.CommunityServiceInterface
	postService      service.PostServiceInterface
	authService      service.AuthServiceInterface
	voteService      service.VoteServiceInterface
	commentService   service.CommentServiceInterface
)

// InitContainer 初始化容器
func InitContainer() {
	once.Do(func() {
		// 初始化 DAO 层
		dao.SetDefault(MySQL.GetDB())

		communityService = service.NewCommunityService(dao.Community.WithContext(context.Background()), dao.NewCommunityDao(MySQL.GetDB()))
		postService = service.NewPostService(dao.Post.WithContext(context.Background()), dao.NewPostDao(MySQL.GetDB()))
		authService = service.NewAuthService(dao.User.WithContext(context.Background()), dao.NewUserDao(MySQL.GetDB()))
		voteService = service.NewVoteService(dao.Vote.WithContext(context.Background()), dao.NewVoteDao(MySQL.GetDB()))
		commentService = service.NewCommentService(dao.NewCommentDao(MySQL.GetDB()))
	})
}

// GetCommunityService 返回接口类型的 Service 实例
func GetCommunityService() service.CommunityServiceInterface {
	if communityService == nil {
		panic("community service is not initialized")
	}
	return communityService
}

// GetPostService 返回接口类型的 Service 实例
func GetPostService() service.PostServiceInterface {
	if postService == nil {
		panic("post service is not initialized")
	}
	return postService
}

func GetAuthService() service.AuthServiceInterface {
	if authService == nil {
		panic("auth service is not initialized")
	}
	return authService
}

func GetVoteService() service.VoteServiceInterface {
	if voteService == nil {
		panic("vote service is not initialized")
	}
	return voteService
}

func GetCommentService() service.CommentServiceInterface {
	if commentService == nil {
		panic("comment service is not initialized")
	}
	return commentService
}
