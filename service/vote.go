package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"context"
	"fmt"
	"time"
)

var _ VoteServiceInterface = (*VoteService)(nil)

type VoteServiceInterface interface {
	Vote(ctx context.Context, id int64, voteFor int, userID int64, voteType int) *apiError.ApiError
	RevokeVote(ctx context.Context, id int64, voteFor int, userID int64) *apiError.ApiError
	MyVoteList(ctx context.Context, userID int64, voteFor, pageNum, pageSize int) ([]int64, *apiError.ApiError)
	GetVoteCount(ctx context.Context, id int64, voteFor int) (int64, int64, *apiError.ApiError)
	CheckUserVoted(ctx context.Context, id []int64, voteFor int, userID int64) ([]model.Vote, *apiError.ApiError)
	GetPostVoteDetail(ctx context.Context, postID int64, pageNum, pageSize int) ([]*DTO.UserVoteDetailDTO, *apiError.ApiError)
	GetCommentVoteDetail(ctx context.Context, commentID int64, pageNum, pageSize int) ([]*DTO.UserVoteDetailDTO, *apiError.ApiError)
}

type VoteService struct {
	dao.IVoteDo
	dao.VoteDaoInterface
}

/*
Vote
投票分为四种情况：1.投赞成票 2.投反对票 3.取消投票 4.反转投票

记录文章参与投票的人
更新文章分数：赞成票要加分；反对票减分

v=1时，有两种情况

	1.之前没投过票，现在要投赞成票
	2.之前投过反对票，现在要改为赞成票

v=0时，有两种情况

	1.之前投过赞成票，现在要取消
	2.之前投过反对票，现在要取消

v=-1时，有两种情况

	1.之前没投过票，现在要投反对票
	2.之前投过赞成票，现在要改为反对票
*/
func (v *VoteService) Vote(ctx context.Context, id int64, voteFor int, userID int64, voteType int) *apiError.ApiError {
	// 查询先前的投票记录
	voteRecord, err := v.VoteDaoInterface.GetVoteRecord(ctx, id, voteFor, userID)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票记录失败",
		}
	}
	caseNum := 0
	if voteRecord == voteType {
		return &apiError.ApiError{
			Code: code.InvalidParam,
			Msg:  "请勿重复投票",
		}
	}

	// 根据投票类型和先前的投票记录，计算投票数的变化量
	if voteType == 1 {
		if voteRecord == 0 {
			err = v.VoteDaoInterface.VoteCase1(ctx, id, voteFor, userID)
			caseNum = 1
		} else {
			err = v.VoteDaoInterface.VoteCase2(ctx, id, voteFor, userID)
			caseNum = 2
		}
	} else if voteType == 0 {
		if voteRecord == 1 {
			err = v.VoteDaoInterface.VoteCase3(ctx, id, voteFor, userID)
			caseNum = 3
		} else {
			err = v.VoteDaoInterface.VoteCase4(id, ctx, voteFor, userID)
			caseNum = 4
		}
	} else {
		if voteRecord == 0 {
			err = v.VoteDaoInterface.VoteCase5(ctx, id, voteFor, userID)
			caseNum = 5
		} else {
			err = v.VoteDaoInterface.VoteCase6(ctx, id, voteFor, userID)
			caseNum = 6
		}
	}

	// 异步更新帖子的投票数
	go updatePostVoteCount(ctx, v, id, voteFor, caseNum)

	return nil
}

// RevokeVote 取消投票
func (v *VoteService) RevokeVote(ctx context.Context, id int64, voteFor int, userID int64) *apiError.ApiError {
	return v.Vote(ctx, id, voteFor, userID, 0)
}

func (v *VoteService) MyVoteList(ctx context.Context, userID int64, voteFor, pageNum, pageSize int) ([]int64, *apiError.ApiError) {
	voteRecord, err := v.VoteDaoInterface.GetUserVoteList(ctx, voteFor, userID, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票记录失败",
		}
	}
	return voteRecord, nil
}

func (v *VoteService) GetVoteCount(ctx context.Context, id int64, voteFor int) (int64, int64, *apiError.ApiError) {
	var up, down int64
	var err error
	if voteFor == dao.VotePost {
		up, down, err = v.VoteDaoInterface.GetContentVoteCount(ctx, id)
	} else {
		up, down, err = v.VoteDaoInterface.GetCommentVoteCount(ctx, id)
	}
	if err != nil {
		return 0, 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票数失败",
		}
	}
	return up, down, nil
}

// CheckUserVoted 批量查询用户是否投票过
func (v *VoteService) CheckUserVoted(ctx context.Context, id []int64, voteFor int, userID int64) ([]model.Vote, *apiError.ApiError) {
	votes, err := v.VoteDaoInterface.CheckUserVoted(ctx, id, voteFor, userID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("批量查询投票记录失败: %v", err),
		}
	}
	return votes, nil
}

func (v *VoteService) GetPostVoteDetail(ctx context.Context, postID int64, pageNum, pageSize int) ([]*DTO.UserVoteDetailDTO, *apiError.ApiError) {
	voteDetails, err := v.VoteDaoInterface.GetPostVoteDetail(ctx, postID, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("查询投票详情失败: %v", err),
		}
	}
	resp := make([]*DTO.UserVoteDetailDTO, len(voteDetails))
	for i, voteDetail := range voteDetails {
		resp[i] = &voteDetail
	}
	return resp, nil
}

func (v *VoteService) GetCommentVoteDetail(ctx context.Context, commentID int64, pageNum, pageSize int) ([]*DTO.UserVoteDetailDTO, *apiError.ApiError) {
	voteDetails, err := v.VoteDaoInterface.GetCommentVoteDetail(ctx, commentID, pageNum, pageSize)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("查询投票详情失败: %v", err),
		}
	}
	resp := make([]*DTO.UserVoteDetailDTO, len(voteDetails))
	for i, voteDetail := range voteDetails {
		resp[i] = &voteDetail
	}
	return resp, nil
}

func NewVoteService(voteDao dao.IVoteDo, voteDaoInterface dao.VoteDaoInterface) VoteServiceInterface {
	return &VoteService{
		voteDao,
		voteDaoInterface,
	}
}

func updatePostVoteCount(ctx context.Context, v *VoteService, id int64, voteFor int, caseNum int) {
	ctx = context.Background()

	var maxRetries = 3          // 最大重试次数
	var delay = 2 * time.Second // 重试间隔时间
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		switch caseNum {
		case 1:
			err = v.VoteDaoInterface.ContentVoteCase1(ctx, id, voteFor)
			break
		case 2:
			err = v.VoteDaoInterface.ContentVoteCase2(ctx, id, voteFor)
			break
		case 3:
			err = v.VoteDaoInterface.ContentVoteCase3(ctx, id, voteFor)
			break
		case 4:
			err = v.VoteDaoInterface.ContentVoteCase4(ctx, id, voteFor)
			break
		case 5:
			err = v.VoteDaoInterface.ContentVoteCase5(ctx, id, voteFor)
			break
		case 6:
			err = v.VoteDaoInterface.ContentVoteCase6(ctx, id, voteFor)
			break
		}
		if err != nil {
			time.Sleep(delay)
			continue
		} else {
			return
		}
	}
}
