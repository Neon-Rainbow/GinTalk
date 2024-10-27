// Package container: 依赖注入容器，用于初始化 Service 实例。
package container

import (
	"GinTalk/cache"
	"GinTalk/controller"
	"GinTalk/dao"
	"GinTalk/dao/MySQL"
	"GinTalk/dao/Redis"
	"GinTalk/kafka"
	"GinTalk/service"

	"github.com/go-redis/redis/v8"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

var container *dig.Container

func BuildContainer() *dig.Container {
	container = dig.New()

	container.Provide(MySQL.GetDB)
	container.Provide(Redis.GetRedisClient)

	// 提供 MySQL 连接
	container.Provide(func(db *gorm.DB) dao.CommentDaoInterface {
		return dao.NewCommentDao(db)
	})
	container.Provide(func(db *gorm.DB) dao.PostDaoInterface {
		return dao.NewPostDao(db)
	})
	container.Provide(func(db *gorm.DB) dao.UserDaoInterface {
		return dao.NewUserDao(db)
	})
	container.Provide(func(db *gorm.DB) dao.CommunityDaoInterface {
		return dao.NewCommunityDao(db)
	})
	container.Provide(func(db *gorm.DB) dao.PostVoteDaoInterface {
		return dao.NewPostVoteDao(db)
	})
	container.Provide(func(db *gorm.DB) dao.CommentVoteInterface {
		return dao.NewCommentVoteImpl(db)
	})

	// 提供 Redis 连接
	container.Provide(func(client *redis.Client) cache.PostCacheInterface {
		return cache.NewPostCache(client)
	})
	container.Provide(func(client *redis.Client) cache.AuthCacheInterface {
		return cache.NewAuthCache(client)
	})
	container.Provide(func(client *redis.Client) cache.VoteCacheInterface {
		return cache.NewVoteCache(client)
	})

	// 提供 Kafka 实例
	container.Provide(func(votePost dao.PostVoteDaoInterface, voteCache cache.VoteCacheInterface, postCache cache.PostCacheInterface) kafka.KafkaInterface {
		return kafka.NewKafka(votePost, voteCache, postCache)
	})

	container.Provide(func(voteDao dao.PostVoteDaoInterface, voteCache cache.VoteCacheInterface) kafka.MessageHandle {
		return kafka.NewVotePostHandle(voteDao, voteCache)
	})

	// 提供 Service 实例
	container.Provide(func(dao dao.CommentDaoInterface) service.CommentServiceInterface {
		return service.NewCommentService(dao)
	})

	container.Provide(func(dao dao.PostDaoInterface, cache cache.PostCacheInterface, kafka kafka.KafkaInterface) service.PostServiceInterface {
		return service.NewPostService(dao, cache, kafka)
	})

	container.Provide(func(dao dao.UserDaoInterface, cache cache.AuthCacheInterface) service.AuthServiceInterface {
		return service.NewAuthService(dao, cache)
	})

	container.Provide(func(dao dao.CommunityDaoInterface) service.CommunityServiceInterface {
		return service.NewCommunityService(dao)
	})

	container.Provide(func(dao dao.PostVoteDaoInterface, cache cache.VoteCacheInterface, kafka kafka.KafkaInterface) service.VotePostServiceInterface {
		return service.NewVoteService(dao, cache, kafka)
	})

	container.Provide(func(dao dao.CommentVoteInterface) service.VoteCommentServiceInterface {
		return service.NewVoteCommentService(dao)
	})

	// 提供 Controller 实例
	container.Provide(func(service service.AuthServiceInterface) *controller.AuthHandler {
		return controller.NewAuthHandler(service)
	})
	container.Provide(func(service service.CommunityServiceInterface) *controller.CommunityHandler {
		return controller.NewCommunityController(service)
	})
	container.Provide(func(service service.PostServiceInterface) *controller.PostHandler {
		return controller.NewPostHandler(service)
	})
	container.Provide(func(service service.VotePostServiceInterface) *controller.VoteHandler {
		return controller.NewVoteHandle(service)
	})
	container.Provide(func(service service.CommentServiceInterface) *controller.CommentController {
		return controller.NewCommentController(service)
	})
	container.Provide(func(service service.VoteCommentServiceInterface) *controller.VoteCommentController {
		return controller.NewVoteCommentController(service)
	})

	return container
}

func GetContainer() *dig.Container {
	return container
}
