package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"context"
	"errors"
	"gorm.io/gorm"
)

// GetCommunityList 获取社区列表
func GetCommunityList(ctx context.Context) ([]*DTO.CommunityNameDTO, *apiError.ApiError) {
	// 使用 DAO 获取社区列表
	communities, err := dao.GetCommunityList(ctx)

	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取社区列表失败",
		}
	}

	// 构造响应数据
	resp := make([]*DTO.CommunityNameDTO, 0)
	for _, community := range communities {
		resp = append(resp, &DTO.CommunityNameDTO{
			CommunityID:   community.CommunityID,
			CommunityName: community.CommunityName,
		})
	}

	return resp, nil
}

// GetCommunityDetail 获取社区详情
func GetCommunityDetail(ctx context.Context, communityID int32) (*DTO.CommunityDetailDTO, *apiError.ApiError) {
	// 使用 DAO 获取社区详情
	community, err := dao.GetCommunityDetail(ctx, communityID)

	// 处理错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &apiError.ApiError{
				Code: code.ServerError,
				Msg:  "社区未找到",
			}
		}
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取社区详情失败",
		}
	}

	return community, nil
}
