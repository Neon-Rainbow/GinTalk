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
	"go.uber.org/zap"
	"slices"
	"time"
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

func GetPostList(ctx context.Context, pageNum int, pageSize int, order int) ([]DTO.PostSummary, *apiError.ApiError) {
	// pageNum 和 pageSize 不能小于等于 0
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

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
