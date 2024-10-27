package controller

import (
	"GinTalk/DTO"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTopComments 获取主评论
// @Summary 获取主评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param post_id query string true "帖子ID"
// @Param page_size query string true "每页数量"
// @Param page_num query string true "页码"
// @Success 200 {object} CommentListResponse
// @Router /api/v1/comment/top [get]
func GetTopComments(c *gin.Context) {
	// 1. 从请求中获取参数
	_postID := c.Query("post_id")
	pageNum, pageSize := getPageInfo(c)

	// 2. 参数校验
	_postIDInt, err := strconv.Atoi(_postID)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "post_id 参数错误")
		return
	}
	postID := int64(_postIDInt)

	// 3. 调用 service 获取数据
	commentList, apiError := service.GetTopComments(c, postID, pageSize, pageNum)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	// 4. 返回响应
	ResponseSuccess(c, commentList)
}

// GetSubComments 获取子评论
// @Summary 获取子评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param post_id query string true "帖子ID"
// @Param parent_id query string true "父评论ID"
// @Param page_size query string true "每页数量"
// @Param page_num query string true "页码"
// @Success 200 {object} CommentListResponse
// @Router /api/v1/comment/sub [get]
func GetSubComments(c *gin.Context) {
	// 1. 从请求中获取参数
	_postID := c.Query("post_id")
	_parentID := c.Query("parent_id")
	pageNum, pageSize := getPageInfo(c)

	// 2. 参数校验
	_postIDInt, err := strconv.Atoi(_postID)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "post_id 参数错误")
		return
	}
	postID := int64(_postIDInt)
	_parentIDInt, err := strconv.Atoi(_parentID)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "parent_id 参数错误")
		return
	}
	parentID := int64(_parentIDInt)
	// 3. 调用 service 获取数据
	commentList, apiError := service.GetSubComments(c, postID, parentID, pageSize, pageNum)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	// 4. 返回响应
	ResponseSuccess(c, commentList)
}

// GetCommentByID 获取评论
// @Summary 获取评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param comment_id query string true "评论ID"
// @Success 200 {object} Comment
// @Router /api/v1/comment [get]
func GetCommentByID(c *gin.Context) {
	// 1. 从请求中获取参数
	_commentID := c.Query("comment_id")

	// 2. 参数校验
	commentID, err := strconv.Atoi(_commentID)
	if err != nil {
		ResponseBadRequest(c, "comment_id 参数错误")
		return
	}
	// 3. 调用 service 获取数据
	comment, apiError := service.GetCommentByID(c, int64(commentID))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	// 4. 返回响应
	ResponseSuccess(c, comment)
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	// 1. 从请求中获取参数
	var comment DTO.CreateCommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "参数错误")
		return
	}
	username, _ := c.Get(ContextUsernameKey)
	comment.AuthorName = username.(string)
	// 2. 参数校验
	if comment.Content == "" {
		ResponseErrorWithMsg(c, code.InvalidParam, "content 参数错误")
		return
	}

	userID, _ := c.Get(ContextUserIDKey)
	comment.AuthorID = userID.(int64)

	// 3. 调用 service 获取数据
	apiError := service.CreateComment(c, &comment)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	// 4. 返回响应
	ResponseSuccess(c, nil)
}

// UpdateComment 更新评论
func UpdateComment(c *gin.Context) {
	// 1. 从请求中获取参数
	var comment DTO.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "参数错误")
		return
	}
	// 2. 参数校验
	if comment.Content == "" {
		ResponseErrorWithMsg(c, code.InvalidParam, "content 参数错误")
		return
	}
	// 3. 调用 service 获取数据
	apiError := service.UpdateComment(c, &comment)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	// 4. 返回响应
	ResponseSuccess(c, nil)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	// 1. 从请求中获取参数
	_commentID := c.Query("comment_id")

	// 2. 参数校验
	commentID, err := strconv.Atoi(_commentID)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "comment_id 参数错误")
		return
	}
	// 3. 调用 service 获取数据
	apiError := service.DeleteComment(c, int64(commentID))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	// 4. 返回响应
	ResponseSuccess(c, nil)
}

func GetCommentCount(c *gin.Context) {
	// 1. 从请求中获取参数
	_postID := c.Query("post_id")

	// 2. 参数校验
	postID, err := strconv.Atoi(_postID)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "post_id 参数错误")
		return
	}
	//3. 调用 service 获取数据
	commentCount, apiError := service.GetCommentCount(c, int64(postID))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	//4. 返回响应
	ResponseSuccess(c, commentCount)
}

func GetTopCommentCount(c *gin.Context) {
	// 1. 从请求中获取参数
	_postID := c.Query("post_id")

	// 2. 参数校验
	postID, err := strconv.Atoi(_postID)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "post_id 参数错误")
		return
	}
	//3. 调用 service 获取数据
	commentCount, apiError := service.GetTopCommentCount(c, int64(postID))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	//4. 返回响应
	ResponseSuccess(c, commentCount)
}

func GetSubCommentCount(c *gin.Context) {
	// 1. 从请求中获取参数
	_parentID := c.Query("parent_id")

	// 2. 参数校验
	parentID, err := strconv.Atoi(_parentID)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "parent_id 参数错误")
		return
	}
	//3. 调用 service 获取数据
	commentCount, apiError := service.GetSubCommentCount(c, int64(parentID))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	//4. 返回响应
	ResponseSuccess(c, commentCount)
}

func GetCommentCountByUserID(c *gin.Context) {
	// 1. 从请求中获取参数
	_userID := c.Query("user_id")

	// 2. 参数校验
	userID, err := strconv.Atoi(_userID)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "user_id 参数错误")
		return
	}
	//3. 调用 service 获取数据
	commentCount, apiError := service.GetCommentCountByUserID(c, int64(userID))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	//4. 返回响应
	ResponseSuccess(c, commentCount)
}

func GetCommentByCommentID(c *gin.Context) {
	// 1. 从请求中获取参数
	_commentID := c.Query("comment_id")

	// 2. 参数校验
	commentID, err := strconv.Atoi(_commentID)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "comment_id 参数错误")
		return
	}
	//3. 调用 service 获取数据
	comment, apiError := service.GetCommentByID(c, int64(commentID))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	//4. 返回响应
	ResponseSuccess(c, comment)
}
