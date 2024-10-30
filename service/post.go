package service

import (
	"GinTalk/DTO"
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/kafka"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"GinTalk/pkg/snowflake"
	"context"
	"fmt"
	"slices"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

// DelayDeleteTime 设置延迟双删的时间
const DelayDeleteTime = 2 * time.Second

func CreatePost(ctx context.Context, postDTO *DTO.PostDetail) *apiError.ApiError {
	postID, err := snowflake.GetID()
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("生成帖子ID失败: %v", err),
		}
	}

	postDTO.PostID = postID

	// 将帖子 ID 存入 Redis
	go func() {
		err := kafka.SendPostMessage(context.Background(), postDTO)
		if err != nil {
			zap.L().Error("Kafka 生产消息失败", zap.Error(err))
		}
	}()

	return nil
}

// GetPostList 根据提供的分页和排序参数检索帖子摘要列表。
// 它使用 singleflight 机制防止缓存雪崩，并尝试首先从 Redis 缓存中获取数据。
// 如果缓存中缺少一些帖子，它会从数据库中获取这些帖子并更新缓存。
//
// 参数:
//   - ctx: 用于管理请求生命周期的上下文。
//   - pageNum: 分页的页码。如果小于或等于 0，则默认为 1。
//   - pageSize: 每页的帖子数量。如果小于或等于 0，则默认为 10。
//   - order: 帖子检索的排序方式。
//
// 返回:
//   - []DTO.PostSummary: 帖子摘要的切片。
//   - *apiError.ApiError: 如果过程中发生错误，则返回错误对象。
func GetPostList(ctx context.Context, pageNum int, pageSize int, order int) ([]DTO.PostSummary, *apiError.ApiError) {
	// pageNum 和 pageSize 不能小于等于 0
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 使用单飞模式, 从 Redis 中获取帖子列表
	sgKey := GenerateSingleFlightKey(SingleFlightKeyPostList, order, pageNum, pageSize)

	// 使用 singleflight 防止缓存雪崩
	var group singleflight.Group

	result, err, _ := group.Do(sgKey, func() (interface{}, error) {
		postIDs, err := cache.GetPostIDs(ctx, order, pageNum, pageSize)
		if err != nil {
			return nil, &apiError.ApiError{
				Code: code.ServerError,
				Msg:  fmt.Sprintf("获取帖子列表失败: %v", err),
			}
		}

		// 首先从 Redis 中获取帖子列表
		redisList, missingIDs, err := cache.GetPostSummary(ctx, postIDs)
		if err != nil {
			return nil, &apiError.ApiError{
				Code: code.ServerError,
				Msg:  fmt.Sprintf("获取帖子列表失败: %v", err),
			}
		}

		// 如果缓存中没有缺失的帖子, 则直接返回
		if len(missingIDs) == 0 {
			return redisList, nil
		}

		list, err := dao.GetPostListBatch(ctx, missingIDs)
		if err != nil {
			return nil, &apiError.ApiError{
				Code: code.ServerError,
				Msg:  fmt.Sprintf("获取帖子列表失败: %v", err),
			}
		}

		// 将缺失的帖子存入 Redis
		go func() {
			for _, post := range list {
				err := cache.SavePost(context.Background(), &post)
				if err != nil {
					zap.L().Error("保存帖子到 Redis 失败", zap.Error(err))
				}
			}
		}()

		return slices.Concat(redisList, list), nil
	})

	if err != nil {
		return nil, err.(*apiError.ApiError)
	}

	return result.([]DTO.PostSummary), nil
}

func GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDetail, *apiError.ApiError) {
	postDetail, err := dao.GetPostDetail(ctx, postID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取帖子详情失败: %v", err),
		}
	}
	return postDetail, nil
}

func UpdatePost(ctx context.Context, postDTO *DTO.PostDetail) *apiError.ApiError {
	// 延迟双删, 保证数据一致性

	// 第一次删除 Redis 中数据
	err := cache.DeletePostSummary(ctx, postDTO.PostID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("删除Redis数据失败, %v", err.Error()),
		}
	}

	if postDTO.PostID == 0 {
		return &apiError.ApiError{
			Code: code.InvalidParam,
			Msg:  "postID 不能为空",
		}
	}

	fmt.Printf("截断前: %s\n", postDTO.Content)
	summary := TruncateByWords(postDTO.Content, MaxSummaryLength)
	fmt.Printf("截断后: %s\n", summary)

	err = dao.UpdatePost(ctx, postDTO, summary)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("更新帖子失败: %v", err),
		}
	}

	// 等待 2s 后第二次删除 Redis 中数据
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), DelayDeleteTime)
		defer cancel()
		time.Sleep(2 * time.Second)
		err := cache.DeletePostSummary(ctx, postDTO.PostID)
		if err != nil {
			zap.L().Error("删除 Redis 数据失败")
		}
	}()

	return nil
}

func GetPostListByCommunityID(ctx context.Context, communityID int64, pageNum int, pageSize int) ([]DTO.PostSummary, *apiError.ApiError) {
	// pageNum 和 pageSize 不能小于等于 0
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	list, err := dao.GetPostListByCommunityID(ctx, communityID, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("获取社区帖子列表失败: %v", err),
		}
	}
	return list, nil
}

func DeletePost(ctx context.Context, postID int64) *apiError.ApiError {
	err := dao.DeletePost(ctx, postID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("删除帖子失败: %v", err),
		}
	}
	go func() {
		err := cache.DeletePost(context.Background(), postID)
		if err != nil {
			zap.L().Error("删除 Redis 中的帖子数据失败, ", zap.Error(err))
		}
	}()

	return nil
}
