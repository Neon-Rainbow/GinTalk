package dao

import (
	"GinTalk/DTO"
	"GinTalk/dao/MySQL"
	"GinTalk/model"
	"context"
)

func GetCommunityList(ctx context.Context) ([]DTO.CommunityListDTO, error) {
	var communities []DTO.CommunityListDTO
	err := MySQL.GetDB().
		WithContext(ctx).
		Model(&model.Community{}).
		Select("community_id, community_name").
		Find(&communities).
		Error
	return communities, err
}

func GetCommunityDetail(ctx context.Context, communityID uint) (*model.Community, error) {
	var community model.Community
	err := MySQL.GetDB().WithContext(ctx).First(&community, communityID).Error
	return &community, err
}
