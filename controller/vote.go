package controller

import (
	"GinTalk/DTO"
	"GinTalk/container"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"github.com/gin-gonic/gin"
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
	var vote DTO.VoteDTO
	if err := c.ShouldBindBodyWithJSON(&vote); err != nil {
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

	apiError := vh.VoteServiceInterface.RevokeVote(c.Request.Context(), vote.PostID, userID)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
	return
}

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

	voteList, apiError := vh.VoteServiceInterface.MyVoteList(c.Request.Context(), myVoteList.UserID, myVoteList.PageNum, myVoteList.PageSize)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}

	ResponseSuccess(c, voteList)
}

func (vh *VoteHandler) GetVoteCountHandler(c *gin.Context) {
	var voteCount DTO.VoteCountDTO
	if err := c.ShouldBindQuery(&voteCount); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	upCount, downCount, apiError := vh.VoteServiceInterface.GetVoteCount(c.Request.Context(), voteCount.PostID)
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

	ResponseSuccess(c, voteList)
	return
}
