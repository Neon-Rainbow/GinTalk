package DTO

// CommunityNameDTO 社区
type CommunityNameDTO struct {
	CommunityID   int32  `json:"community_id"`
	CommunityName string `json:"community_name"`
}

// CommunityDetailDTO 社区详情
type CommunityDetailDTO struct {
	*CommunityNameDTO
	Introduction string `json:"introduction"`
}
