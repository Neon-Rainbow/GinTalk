package controller

import (
	"GinTalk/DTO"
	"GinTalk/container"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type PostHandler struct {
	service.PostServiceInterface
}

func NewPostHandler() *PostHandler {
	return &PostHandler{
		container.GetPostService(),
	}
}

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
func (ph *PostHandler) CreatePostHandler(c *gin.Context) {
	var post DTO.PostDetail
	if err := c.ShouldBindJSON(&post); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	if !isUserIDMatch(c, post.AuthorId) {
		ResponseErrorWithMsg(c, code.InvalidAuth, "无权限操作")
		return
	}

	apiError := ph.PostServiceInterface.CreatePost(c.Request.Context(), &post)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
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
func (ph *PostHandler) GetPostListHandler(c *gin.Context) {
	pageNum, pageSize := getPageInfo(c)
	postList, apiError := ph.PostServiceInterface.GetPostList(c.Request.Context(), pageNum, pageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
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
func (ph *PostHandler) GetPostDetailHandler(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("ID"), 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	post, apiError := ph.PostServiceInterface.GetPostDetail(c.Request.Context(), postID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
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
func (ph *PostHandler) UpdatePostHandler(c *gin.Context) {
	var post DTO.PostDetail
	if err := c.ShouldBindJSON(&post); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	if !isUserIDMatch(c, post.AuthorId) {
		ResponseErrorWithMsg(c, code.InvalidAuth, "无权限操作")
		return
	}
	apiError := ph.UpdatePost(c.Request.Context(), &post)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
}

func getPageInfo(c *gin.Context) (pageNum int, pageSize int) {
	var err error
	_n, _s := c.Query("page_num"), c.Query("page_size")
	pageNum, err = strconv.Atoi(_n)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}
	pageSize, err = strconv.Atoi(_s)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	return pageNum, pageSize
}
