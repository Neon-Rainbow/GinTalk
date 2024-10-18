package VO

import "GinTalk/model"

// CommunityVO 社区
type CommunityVO struct {
	CommunityID   int32  `json:"community_id"`
	CommunityName string `json:"community_name"`
}

// CommunityDetailVO 社区详情
type CommunityDetailVO struct {
	*model.Community
}
