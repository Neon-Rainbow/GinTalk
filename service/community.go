package service

import (
	"GinTalk/VO"
	"GinTalk/dao"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"context"
	"errors"
	"gorm.io/gorm"
)

// CommunityServiceInterface 的实现
var _ CommunityServiceInterface = (*CommunityService)(nil)

// CommunityServiceInterface 定义社区服务的接口
type CommunityServiceInterface interface {
	GetCommunityList(ctx context.Context) ([]*VO.CommunityVO, *apiError.ApiError)
	GetCommunityDetail(ctx context.Context, communityID int32) (*VO.CommunityDetailVO, *apiError.ApiError)
}

// CommunityService 是 CommunityServiceInterface 的实现
type CommunityService struct {
	communityDao dao.ICommunityDo
}

// NewCommunityService 使用依赖注入初始化 CommunityService
func NewCommunityService(communityDao dao.ICommunityDo) CommunityServiceInterface {
	return &CommunityService{
		communityDao: communityDao,
	}
}

// GetCommunityList 获取社区列表
func (s *CommunityService) GetCommunityList(ctx context.Context) ([]*VO.CommunityVO, *apiError.ApiError) {
	// 使用 DAO 获取社区列表
	communities, err := s.communityDao.WithContext(ctx).
		Select(dao.Community.CommunityID, dao.Community.CommunityName).
		Find()

	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "获取社区列表失败",
		}
	}

	// 构造响应数据
	resp := make([]*VO.CommunityVO, 0)
	for _, community := range communities {
		resp = append(resp, &VO.CommunityVO{
			CommunityID:   community.CommunityID,
			CommunityName: community.CommunityName,
		})
	}

	return resp, nil
}

// GetCommunityDetail 获取社区详情
func (s *CommunityService) GetCommunityDetail(ctx context.Context, communityID int32) (*VO.CommunityDetailVO, *apiError.ApiError) {
	// 使用 DAO 获取社区详情
	community, err := s.communityDao.WithContext(ctx).
		Where(dao.Community.CommunityID.Eq(communityID)).
		First()

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

	// 构造响应数据
	resp := &VO.CommunityDetailVO{
		Community: community,
	}

	return resp, nil
}
