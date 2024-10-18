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
)

// InitContainer 初始化容器
func InitContainer() {
	once.Do(func() {
		// 获取数据库实例
		dbInstance := MySQL.GetDB()

		// 初始化 DAO 层
		dao.SetDefault(dbInstance)

		communityService = service.NewCommunityService(dao.Community.WithContext(context.Background()))
		postService = service.NewPostService(dao.Post.WithContext(context.Background()))
		authService = service.NewAuthService(dao.User.WithContext(context.Background()))
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
