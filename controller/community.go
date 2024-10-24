package controller

import (
	"GinTalk/pkg/code"
	"GinTalk/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 社区控制器
type CommunityHandler struct {
	CommunityService service.CommunityServiceInterface
}

// NewCommunityController 创建 CommunityHandler 实例
func NewCommunityController(service service.CommunityServiceInterface) *CommunityHandler {
	return &CommunityHandler{
		CommunityService: service,
	}
}

func (cc *CommunityHandler) CommunityHandler(c *gin.Context) {
	list, apiError := cc.CommunityService.GetCommunityList(c.Request.Context())
	if apiError != nil {
		zap.L().Error("service.GetCommunityList(c.Request.Context()) 错误",
			zap.Error(apiError),
		)
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, list)
}

func (cc *CommunityHandler) CommunityDetailHandler(c *gin.Context) {
	_s := c.Param("id")

	//string 转为 int32
	_t, err := strconv.Atoi(_s)
	if err != nil {
		zap.L().Error("strconv.Atoi(_s) 错误", zap.Error(err))
		ResponseErrorWithCode(c, code.InvalidParam)
		return
	}

	communityID := int32(_t)

	community, apiError := cc.CommunityService.GetCommunityDetail(c.Request.Context(), communityID)
	if apiError != nil {
		zap.L().Error("service.GetCommunityDetail(c.Request.Context(), communityID) 错误",
			zap.Error(apiError),
		)
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, community)
	return
}
