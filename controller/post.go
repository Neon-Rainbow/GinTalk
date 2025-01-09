package controller

import (
	"GinTalk/DTO"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"go.uber.org/zap"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建帖子
// @Tags 帖子
// @Accept json
// @Produce json
// @Param Authorization header string true "
// @Param post body DTO.PostDetail true "帖子信息"
// @Success 200 {object} Response
// @Router /api/v1/post [post]
func CreatePostHandler(c *gin.Context) {
	var post DTO.PostDetail
	if err := c.ShouldBindJSON(&post); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("CreatePostHandler.ShouldBindJSON() 失败", zap.Error(err))
		return
	}
	if !isUserIDMatch(c, post.AuthorId) {
		ResponseErrorWithMsg(c, code.InvalidAuth, "无权限操作")
		zap.L().Info("CreatePostHandler.isUserIDMatch() 失败")
		return
	}

	if apiError := service.CreatePost(c.Request.Context(), &post); apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("kafka.SendPostMessage() 失败", zap.Error(apiError))
		return
	}
	ResponseSuccess(c, nil)
}

// GetPostListHandler 获取帖子列表
// @Summary 获取帖子列表
// @Description 获取帖子列表
// @Tags 帖子
// @Accept json
// @Produce json
// @Param Authorization header string true "
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} Response
// @Router /api/v1/post [get]
func GetPostListHandler(c *gin.Context) {
	pageNum, pageSize := getPageInfo(c)
	order, err := strconv.Atoi(c.Query("order"))
	if err != nil {
		ResponseBadRequest(c, "order 字段不正确")
		zap.L().Info("GetPostListHandler strconv.Atoi() 失败", zap.Error(err))
		return
	}
	postList, apiError := service.GetPostList(c.Request.Context(), pageNum, pageSize, order)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("PostServiceInterface.GetPostList() 失败", zap.Error(apiError))
		return
	}
	ResponseSuccess(c, postList)
	return
}

// GetPostListByCommunityID 根据社区ID获取帖子列表
// @Summary 根据社区ID获取帖子列表
// @Description 根据社区ID获取帖子列表
// @Tags 帖子
// @Accept json
// @Produce json
// @Param Authorization header string true "
// @Param community_id query int true "社区ID"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} Response
// @Router /api/v1/post/community [get]
func GetPostListByCommunityID(c *gin.Context) {
	pageNum, pageSize := getPageInfo(c)
	communityID, err := strconv.ParseInt(c.Query("community_id"), 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Info("GetPostListByCommunityID strconv.ParseInt() 失败", zap.Error(err))
		return
	}
	postList, apiError := service.GetPostListByCommunityID(c.Request.Context(), communityID, pageNum, pageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("PostServiceInterface.GetPostListByCommunityID() 失败", zap.Error(apiError))
		return
	}
	ResponseSuccess(c, postList)
}

// GetPostDetailHandler 获取帖子详情
// @Summary 获取帖子详情
// @Description 获取帖子详情
// @Tags 帖子
// @Accept json
// @Produce json
// @Param Authorization header string true "
// @Param ID path int true "帖子ID"
// @Success 200 {object} Response
// @Router /api/v1/post/{ID} [get]
func GetPostDetailHandler(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Info("GetPostDetailHandler strconv.ParseInt() 失败", zap.Error(err))
		return
	}

	post, apiError := service.GetPostDetail(c.Request.Context(), postID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("PostServiceInterface.GetPostDetail() 失败", zap.Error(apiError))
		return
	}
	ResponseSuccess(c, post)
}

// UpdatePostHandler 更新帖子
// @Summary 更新帖子
// @Description 更新帖子
// @Tags 帖子
// @Accept json
// @Produce json
// @Param Authorization header string true "
// @Param post body DTO.PostDetail true "帖子信息"
// @Success 200 {object} Response
// @Router /api/v1/post [put]
func UpdatePostHandler(c *gin.Context) {
	var post DTO.PostDetail
	if err := c.ShouldBindJSON(&post); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("UpdatePostHandler.ShouldBindJSON() 失败", zap.Error(err))
		return
	}
	if !isUserIDMatch(c, post.AuthorId) {
		ResponseErrorWithMsg(c, code.InvalidAuth, "无权限操作")
		zap.L().Info("UpdatePostHandler.isUserIDMatch() 失败")
		return
	}
	if apiError := service.UpdatePost(c.Request.Context(), &post); apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("PostServiceInterface.UpdatePost() 失败", zap.Error(apiError))
		return
	}
	ResponseSuccess(c, nil)
}

func DeletePostHandler(c *gin.Context) {
	var p DTO.VotePostDTO
	if err := c.ShouldBindBodyWithJSON(&p); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("DeletePostHandler.ShouldBindBodyWithJSON() 失败", zap.Error(err))
		return
	}
	if !isUserIDMatch(c, p.UserID) {
		ResponseErrorWithMsg(c, code.InvalidAuth, "无权限操作")
		zap.L().Info("DeletePostHandler.isUserIDMatch() 失败")
		return
	}

	if apiError := service.DeletePost(c.Request.Context(), p.PostID); apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("PostServiceInterface.DeletePost() 失败", zap.Error(apiError))
		return
	}
	ResponseSuccess(c, nil)
}

// getPageInfo 获取分页信息
func getPageInfo(c *gin.Context) (pageNum int, pageSize int) {
	var err error
	_n := c.Query("page_num")
	_s := c.Query("page_size")
	if _n == "" {
		_n = c.Query("pageNum")
	}
	if _s == "" {
		_s = c.Query("pageSize")
	}
	pageNum, err = strconv.Atoi(_n)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}
	pageSize, err = strconv.Atoi(_s)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 20
	}
	return pageNum, pageSize
}
