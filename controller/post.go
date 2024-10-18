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

func (ph *PostHandler) CreatePostHandler(c *gin.Context) {
	var post DTO.PostDetail
	if err := c.ShouldBindJSON(&post); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidAuth, err.Error())
		return
	}

	post.AuthorId = userID

	apiError := ph.PostServiceInterface.CreatePost(c.Request.Context(), &post)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
}

func (ph *PostHandler) GetPostListHandler(c *gin.Context) {
	pageNum, pageSize := getPageInfo(c)
	postList, apiError := ph.PostServiceInterface.GetPostList(c.Request.Context(), pageNum, pageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, postList)
}

func (ph *PostHandler) GetPostDetailHandler(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

// GetPostDetailHandler 获取帖子详情
func getPageInfo(c *gin.Context) (pageNum int, pageSize int) {
	var err error
	_n, _s := c.Query("pageNum"), c.Query("pageSize")
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
