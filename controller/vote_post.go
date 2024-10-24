package controller

import (
	"GinTalk/DTO"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VoteHandler struct {
	service.VotePostServiceInterface
}

func NewVoteHandle(service service.VotePostServiceInterface) *VoteHandler {
	return &VoteHandler{
		VotePostServiceInterface: service,
	}
}

// VotePostHandler 投票
// @Summary 投票
// @Description 投票
// @Tags 投票
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param post body DTO.VotePostDTO true "投票信息"
// @Success 200 {object} Response
// @Router /vote/post [post]
func (vh *VoteHandler) VotePostHandler(c *gin.Context) {
	var vote DTO.VotePostDTO
	if err := c.ShouldBindJSON(&vote); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	apiError := vh.VotePostServiceInterface.VotePost(c.Request.Context(), vote.PostID, vote.UserID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
}

// RevokeVoteHandler 取消投票
// @Summary 取消投票
// @Description 取消投票
// @Tags 投票
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param post body DTO.VotePostDTO true "投票信息"
// @Success 200 {object} Response
// @Router /vote/post [delete]
func (vh *VoteHandler) RevokeVoteHandler(c *gin.Context) {
	var vote DTO.VotePostDTO
	if err := c.ShouldBindJSON(&vote); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	apiError := vh.VotePostServiceInterface.RevokeVotePost(c.Request.Context(), vote.PostID, vote.UserID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
}

// MyVoteListHandler 我的投票列表
// @Summary 我的投票列表
// @Description 我的投票列表
// @Tags 投票
// @Accept json
// @Produce json
// @Param user_id query int true "用户 ID"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} Response
// @Router /vote/post/list [get]
func (vh *VoteHandler) MyVoteListHandler(c *gin.Context) {
	userIDInt, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	userID := int64(userIDInt)

	if !isUserIDMatch(c, userID) {
		ResponseErrorWithMsg(c, code.InvalidParam, "用户 ID 不匹配")
		return
	}
	pageNum, pageSize := getPageInfo(c)

	voteList, apiError := vh.VotePostServiceInterface.MyVotePostList(c.Request.Context(), userID, pageNum, pageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, voteList)
}

// GetVoteCountHandler 获取投票数
// @Summary 获取投票数
// @Description 获取投票数
// @Tags 投票
// @Accept json
// @Produce json
// @Param ID path int true "ID"
// @Param vote_for query int true "投票类型"
// @Success 200 {object} Response
// @Router /vote/post/{ID} [get]
func (vh *VoteHandler) GetVoteCountHandler(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	upCount, apiError := vh.VotePostServiceInterface.GetVotePostCount(c.Request.Context(), int64(postID))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, upCount)
}

// CheckUserVotedHandler 检查用户是否投票
// @Summary 检查用户是否投票
// @Description 检查用户是否投票
// @Tags 投票
// @Accept json
// @Produce json
// @Param ID query int true "帖子 ID"
// @Param user_id query int true "用户 ID"
// @Param vote_for query int true "投票类型"
// @Success 200 {object} Response
// @Router /vote/post/list [get]
func (vh *VoteHandler) CheckUserVotedHandler(c *gin.Context) {
	var postIds DTO.CheckVoteListDTO
	if err := c.ShouldBindQuery(&postIds); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	if !isUserIDMatch(c, postIds.UserID) {
		ResponseErrorWithMsg(c, code.InvalidParam, "用户 ID 不匹配")
		return
	}

	voteList, apiError := vh.VotePostServiceInterface.CheckUserPostVoted(c.Request.Context(), postIds.IDs, postIds.UserID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, voteList)
	return
}

// GetPostVoteDetailHandler 获取帖子投票详情
// @Summary 获取帖子投票详情
// @Description 获取帖子投票详情
// @Tags 投票
// @Accept json
// @Produce json
// @Param ID query int true "帖子 ID"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} Response
// @Router /vote/post/detail [get]
func (vh *VoteHandler) GetPostVoteDetailHandler(c *gin.Context) {
	type PostID struct {
		ID int64 `form:"id" binding:"required"`
	}
	pageNum, pageSize := getPageInfo(c)
	var postIds PostID
	if err := c.ShouldBindQuery(&postIds); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	voteList, apiError := vh.VotePostServiceInterface.GetPostVoteDetail(c.Request.Context(), postIds.ID, pageNum, pageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, voteList)
	return
}

func (vh *VoteHandler) GetBatchPostVoteCount(c *gin.Context) {
	type ids struct {
		PostID []int64 `json:"post_id"`
	}
	var postIDs ids
	if err := c.ShouldBindQuery(postIDs); err != nil {
		ResponseErrorWithCode(c, code.InvalidParam)
		return
	}
	resp, apiError := vh.VotePostServiceInterface.GetBatchPostVoteCount(c.Request.Context(), postIDs.PostID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, resp)
	return
}
