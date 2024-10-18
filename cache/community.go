package cache

import (
	"GinTalk/DTO"
	"GinTalk/utils"
	"context"
	"encoding/json"
)

// GetCommunityList 保存社区列表到 redis
func GetCommunityList(ctx context.Context) ([]DTO.CommunityListDTO, error) {
	communityListData, err := getKeyValue[string](ctx, utils.GenerateRedisKey(utils.CommunityListTemplate))
	if err != nil {
		return nil, err
	}

	var communities []DTO.CommunityListDTO
	if err = json.Unmarshal([]byte(communityListData), &communities); err != nil {
		return nil, err
	}
	return communities, nil

}
