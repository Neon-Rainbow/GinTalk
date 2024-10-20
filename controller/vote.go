package controller

import (
	"GinTalk/DTO"
	"GinTalk/container"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type VoteHandler struct {
	service.VoteServiceInterface
}

func NewVoteHandle() *VoteHandler {
	return &VoteHandler{
		container.GetVoteService(),
	}
}

// VoteHandler 投票
// @Summary 投票
// @Description 投票
// @Tags 投票
// @Accept json
// @Produce json
// @Param vote body VoteDTO true "投票"
// @Success 200 {object} Response
// @Router /vote [post]
func (vh *VoteHandler) VoteHandler(c *gin.Context) {
	var vote DTO.VoteDTO
	if err := c.ShouldBindJSON(&vote); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidAuth, err.Error())
		return
	}

	if userID != vote.UserID {
		ResponseErrorWithMsg(c, code.InvalidParam, "用户 ID 不匹配")
		return
	}

	apiError := vh.VoteServiceInterface.Vote(c.Request.Context(), vote.ID, vote.Vote, userID, vote.Vote)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, nil)
	return
}

func (vh *VoteHandler) RevokeVoteHandler(c *gin.Context) {
	// Vote 函数中已经处理了撤销投票的请求,因此不需要再使用该接口
	vh.VoteHandler(c)
	return
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
// @Router /vote/list [get]
func (vh *VoteHandler) MyVoteListHandler(c *gin.Context) {
	var myVoteList DTO.MyVoteListDTO
	if err := c.ShouldBindQuery(&myVoteList); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidAuth, err.Error())
		return
	}

	if userID != myVoteList.UserID {
		ResponseErrorWithMsg(c, code.InvalidParam, "用户 ID 不匹配")
		return
	}
	if myVoteList.PageNum == 0 {
		myVoteList.PageNum = 1
	}
	if myVoteList.PageSize == 0 {
		myVoteList.PageSize = 10
	}

	voteList, apiError := vh.VoteServiceInterface.MyVoteList(
		c.Request.Context(),
		myVoteList.UserID,
		myVoteList.VoteFor,
		myVoteList.PageNum,
		myVoteList.PageSize)
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
// @Router /vote/{ID} [get]
func (vh *VoteHandler) GetVoteCountHandler(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	voteFor, err := strconv.Atoi(c.Query("vote_for"))
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	upCount, downCount, apiError := vh.VoteServiceInterface.GetVoteCount(c.Request.Context(), int64(postID), voteFor)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, gin.H{
		"up_count":   upCount,
		"down_count": downCount,
	})
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
// @Router /vote/list [get]
func (vh *VoteHandler) CheckUserVotedHandler(c *gin.Context) {
	var postIds DTO.CheckVoteListDTO
	if err := c.ShouldBindQuery(&postIds); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidAuth, err.Error())
		return
	}

	if userID != postIds.UserID {
		ResponseErrorWithMsg(c, code.InvalidParam, "用户 ID 不匹配")
		return
	}

	voteList, apiError := vh.VoteServiceInterface.CheckUserVoted(
		c.Request.Context(),
		postIds.IDs,
		postIds.VoteFor,
		userID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	var respList []gin.H

	for _, vote := range voteList {
		respList = append(respList, gin.H{
			"ID":   vote.PostID + vote.CommentID,
			"vote": vote.Vote,
		})
	}

	ResponseSuccess(c, respList)
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
		ID       int64 `form:"id" binding:"required"`
		PageNum  int   `form:"page_num"`
		PageSize int   `form:"page_size"`
	}
	var postIds PostID
	if err := c.ShouldBindQuery(&postIds); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	if postIds.PageNum == 0 {
		postIds.PageNum = 1
	}

	if postIds.PageSize == 0 {
		postIds.PageSize = 10
	}

	voteList, apiError := vh.VoteServiceInterface.GetPostVoteDetail(c.Request.Context(), postIds.ID, postIds.PageNum, postIds.PageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, voteList)
	return
}

// GetCommentVoteDetailHandler 获取评论投票详情
// @Summary 获取评论投票详情
// @Description 获取评论投票详情
// @Tags 投票
// @Accept json
// @Produce json
// @Param ID query int true "评论 ID"
// @Param page_num query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} Response
// @Router /vote/comment/detail [get]
func (vh *VoteHandler) GetCommentVoteDetailHandler(c *gin.Context) {
	type CommentID struct {
		ID       int64 `form:"id" binding:"required"`
		PageNum  int   `form:"page_num"`
		PageSize int   `form:"page_size"`
	}
	var commentIds CommentID
	if err := c.ShouldBindQuery(&commentIds); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	if commentIds.PageNum == 0 {
		commentIds.PageNum = 1
	}

	if commentIds.PageSize == 0 {
		commentIds.PageSize = 10
	}

	voteList, apiError := vh.VoteServiceInterface.GetCommentVoteDetail(c.Request.Context(), commentIds.ID, commentIds.PageNum, commentIds.PageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, voteList)
	return
}
