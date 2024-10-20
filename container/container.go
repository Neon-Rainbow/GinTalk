// Package container: 依赖注入容器，用于初始化 Service 实例。
package container

import (
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/dao/MySQL"
	"GinTalk/dao/Redis"
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
		postService = service.NewPostService(dao.NewPostDao(MySQL.GetDB()), cache.NewPostCache(Redis.GetRedisClient()))
		authService = service.NewAuthService(dao.NewUserDao(MySQL.GetDB()), cache.NewAuthCache(Redis.GetRedisClient()))
		voteService = service.NewVoteService(dao.NewVoteDao(MySQL.GetDB()), cache.NewVoteCache(Redis.GetRedisClient()))
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
