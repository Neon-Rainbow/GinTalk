package controller

import (
	"GinTalk/DTO"
	"GinTalk/pkg/code"
	"GinTalk/service"

	"github.com/gin-gonic/gin"
)

// VoteCommentController 投票评论
// @Summary 投票评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param voteComment body DTO.VoteComment true "voteComment"
// @Success 200 {object} Response
// @Router /vote/comment [post]
func VoteCommentController(c *gin.Context) {
	var voteComment DTO.VoteComment
	if err := c.ShouldBindJSON(&voteComment); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	if apiError := service.VoteComment(voteComment.UserID, voteComment.CommentID); apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
}

// RemoveVoteCommentController 取消评论投票
// @Summary 取消评论投票
// @Tags 评论
// @Accept json
// @Produce json
// @Param voteComment body DTO.VoteComment true "voteComment"
// @Success 200 {object} Response
// @Router /vote/comment [delete]
func RemoveVoteCommentController(c *gin.Context) {
	var voteComment DTO.VoteComment
	if err := c.ShouldBindJSON(&voteComment); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	if apiError := service.RemoveVoteComment(voteComment.UserID, voteComment.CommentID); apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
}

// GetVoteCommentController 获取评论投票数
// @Summary 获取评论投票数
// @Tags 评论
// @Accept json
// @Produce json
// @Param voteComment body DTO.VoteComment true "voteComment"
// @Success 200 {object} Response
// @Router /vote/comment [get]
func GetVoteCommentController(c *gin.Context) {
	var voteComment DTO.VoteComment
	if err := c.ShouldBindJSON(&voteComment); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	count, apiError := service.GetVoteComment(voteComment.UserID, voteComment.CommentID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, count)
}

// GetVoteCommentListController 获取评论投票状态
// @Summary 获取评论投票状态
// @Tags 评论
// @Accept json
// @Produce json
// @Param voteComment body VoteCommentList true "voteComment"
// @Success 200 {object} Response
// @Router /vote/comment/list [get]
func GetVoteCommentListController(c *gin.Context) {
	type VoteCommentList struct {
		UserID    int64   `json:"user_id" form:"user_id"`
		CommentID []int64 `json:"comment_id" form:"comment_id"`
	}
	var voteComment VoteCommentList
	if err := c.ShouldBindQuery(&voteComment); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	result, apiError := service.GetVoteCommentList(voteComment.UserID, voteComment.CommentID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, result)
}
