package dao

import (
	"GinTalk/DTO"
	"GinTalk/dao/MySQL"
	"context"
)

func GetCommunityList(ctx context.Context) ([]*DTO.CommunityNameDTO, error) {
	var communities []*DTO.CommunityNameDTO
	sqlStr := `SELECT community_id, community_name FROM community`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr).Scan(&communities).Error
	if err != nil {
		return nil, err
	}
	return communities, nil
}

func GetCommunityDetail(ctx context.Context, communityID int32) (*DTO.CommunityDetailDTO, error) {
	var communityDetail DTO.CommunityDetailDTO
	sqlStr := `SELECT community_id, community_name, introduction FROM community WHERE community_id = ? AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, communityID).Scan(&communityDetail).Error
	if err != nil {
		return nil, err
	}
	return &communityDetail, nil
}
