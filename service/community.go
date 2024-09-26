package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"context"
)

func GetCommunityList(ctx context.Context) ([]DTO.CommunityListDTO, *apiError.ApiError) {
	communities, err := dao.GetCommunityList(ctx)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取社区列表失败",
		}
	}
	return communities, nil
}

func GetCommunityDetail(ctx context.Context, communityID uint) (*model.Community, *apiError.ApiError) {
	community, err := dao.GetCommunityDetail(ctx, communityID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取社区详情失败",
		}
	}
	return community, nil
}
