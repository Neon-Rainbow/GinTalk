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

func NewVoteHandler() *VoteHandler {
	return &VoteHandler{
		container.GetVoteService(),
	}
}

// VoteHandler 投票
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

	apiError := vh.VoteServiceInterface.Vote(c.Request.Context(), vote.PostID, userID, vote.Vote)
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

func (vh *VoteHandler) MyVoteListHandler(c *gin.Context) {
	var myVoteList DTO.MyVoteListDTO
	if err := c.ShouldBindBodyWithJSON(&myVoteList); err != nil {
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

	voteList, apiError := vh.VoteServiceInterface.MyVoteList(c.Request.Context(), myVoteList.UserID, myVoteList.PageNum, myVoteList.PageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, voteList)
}

func (vh *VoteHandler) GetVoteCountHandler(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	upCount, downCount, apiError := vh.VoteServiceInterface.GetVoteCount(c.Request.Context(), int64(postID))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, gin.H{
		"up_count":   upCount,
		"down_count": downCount,
	})
}

func (vh *VoteHandler) CheckUserVotedHandler(c *gin.Context) {
	var postIds DTO.CreateVoteDTO
	if err := c.ShouldBindJSON(&postIds); err != nil {
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

	voteList, apiError := vh.VoteServiceInterface.CheckUserVoted(c.Request.Context(), postIds.PostID, userID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	type resp struct {
		PostID int64 `json:"post_id"`
		Vote   int   `json:"vote"`
	}
	var respList []gin.H

	for _, vote := range voteList {
		respList = append(respList, gin.H{
			"post_id": vote.PostID,
			"vote":    vote.Vote,
		})
	}

	ResponseSuccess(c, respList)
	return
}

func (vh *VoteHandler) GetPostVoteDetailHandler(c *gin.Context) {
	type PostID struct {
		PostID   int64 `form:"post_id" binding:"required"`
		pageNum  int   `form:"page_num"`
		pageSize int   `form:"page_size"`
	}
	var postIds PostID
	if err := c.ShouldBindQuery(postIds); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	if postIds.pageNum == 0 {
		postIds.pageNum = 1
	}

	if postIds.pageSize == 0 {
		postIds.pageSize = 10
	}

	voteList, apiError := vh.VoteServiceInterface.GetPostVoteDetail(c.Request.Context(), postIds.PostID, postIds.pageNum, postIds.pageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, voteList)
	return
}
