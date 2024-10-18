package controller

import (
	"GinTalk/DTO"
	"GinTalk/container"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"github.com/gin-gonic/gin"
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
	var post DTO.PostDTO
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
