package dao

import (
	"GinTalk/VO"
	"context"
	"gorm.io/gorm"
)

var _ CommunityDaoInterface = (*CommunityDao)(nil)

type CommunityDaoInterface interface {
	GetCommunityList(ctx context.Context) ([]*VO.CommunityVO, error)
	GetCommunityDetail(ctx context.Context, communityID int32) (*VO.CommunityDetailVO, error)
}

type CommunityDao struct {
	*gorm.DB
}

func NewCommunityDao(db *gorm.DB) CommunityDaoInterface {
	return &CommunityDao{DB: db}
}

func (cd *CommunityDao) GetCommunityList(ctx context.Context) ([]*VO.CommunityVO, error) {
	var communities []*VO.CommunityVO
	sqlStr := `SELECT community_id, community_name FROM community`
	err := cd.WithContext(ctx).Raw(sqlStr).Scan(&communities).Error
	if err != nil {
		return nil, err
	}
	return communities, nil
}

func (cd *CommunityDao) GetCommunityDetail(ctx context.Context, communityID int32) (*VO.CommunityDetailVO, error) {
	var communityDetail VO.CommunityDetailVO
	sqlStr := `SELECT community_id, community_name, introduction FROM community WHERE community_id = ? AND delete_time = 0`
	err := cd.WithContext(ctx).Raw(sqlStr, communityID).Scan(&communityDetail).Error
	if err != nil {
		return nil, err
	}
	return &communityDetail, nil
}
