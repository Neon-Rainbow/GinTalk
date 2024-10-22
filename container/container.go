// Package container: 依赖注入容器，用于初始化 Service 实例。
package container

import (
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/dao/MySQL"
	"GinTalk/dao/Redis"
	"GinTalk/service"
	"sync"
)

var (
	once               sync.Once
	communityService   service.CommunityServiceInterface
	postService        service.PostServiceInterface
	authService        service.AuthServiceInterface
	votePostService    service.VotePostServiceInterface
	commentService     service.CommentServiceInterface
	voteCommentService service.VoteCommentServiceInterface
)

// InitContainer 初始化容器
func InitContainer() {
	once.Do(func() {
		// 初始化 DAO 层
		dao.SetDefault(MySQL.GetDB())

		communityService = service.NewCommunityService(dao.NewCommunityDao(MySQL.GetDB()))
		postService = service.NewPostService(dao.NewPostDao(MySQL.GetDB()), cache.NewPostCache(Redis.GetRedisClient()))
		authService = service.NewAuthService(dao.NewUserDao(MySQL.GetDB()), cache.NewAuthCache(Redis.GetRedisClient()))
		votePostService = service.NewVoteService(dao.NewPostVoteDao(MySQL.GetDB()), cache.NewVoteCache(Redis.GetRedisClient()))
		commentService = service.NewCommentService(dao.NewCommentDao(MySQL.GetDB()))
		voteCommentService = service.NewVoteCommentService(dao.NewCommentVoteImpl(MySQL.GetDB()))
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

func GetVotePostService() service.VotePostServiceInterface {
	if votePostService == nil {
		panic("vote service is not initialized")
	}
	return votePostService
}

func GetCommentService() service.CommentServiceInterface {
	if commentService == nil {
		panic("comment service is not initialized")
	}
	return commentService
}

func GetVoteCommentService() service.VoteCommentServiceInterface {
	if voteCommentService == nil {
		panic("vote comment service is not initialized")
	}
	return voteCommentService
}
