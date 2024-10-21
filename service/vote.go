package service

import (
	"GinTalk/DTO"
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"time"
)

var _ VoteServiceInterface = (*VoteService)(nil)

type VoteServiceInterface interface {
	Vote(ctx context.Context, id int64, voteFor int, userID int64, voteType int) *apiError.ApiError
	RevokeVote(ctx context.Context, id int64, voteFor int, userID int64) *apiError.ApiError
	MyVoteList(ctx context.Context, userID int64, voteFor, pageNum, pageSize int) ([]int64, *apiError.ApiError)
	GetVoteCount(ctx context.Context, id int64, voteFor int) (int64, int64, *apiError.ApiError)

	// GetBatchPostVoteCount 该函数用于批量查询帖子的投票数量
	GetBatchPostVoteCount(ctx context.Context, ids []int64) ([]DTO.PostVotes, *apiError.ApiError)
	CheckUserVoted(ctx context.Context, id []int64, voteFor int, userID int64) ([]model.Vote, *apiError.ApiError)
	GetPostVoteDetail(ctx context.Context, postID int64, pageNum, pageSize int) ([]*DTO.UserVoteDetailDTO, *apiError.ApiError)
	GetCommentVoteDetail(ctx context.Context, commentID int64, pageNum, pageSize int) ([]*DTO.UserVoteDetailDTO, *apiError.ApiError)
}

type VoteService struct {
	dao.VoteDaoInterface
	cache.VoteCacheInterface
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

func (v *VoteService) GetBatchPostVoteCount(ctx context.Context, ids []int64) ([]DTO.PostVotes, *apiError.ApiError) {
	resp, err := v.VoteDaoInterface.GetBatchPostVoteCount(ctx, ids)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询错误",
		}
	}
	return resp, nil
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

func NewVoteService(voteDaoInterface dao.VoteDaoInterface, voteCacheInterface cache.VoteCacheInterface) VoteServiceInterface {
	return &VoteService{
		voteDaoInterface,
		voteCacheInterface,
	}
}

const (
	MaxRetries   = 3               // 最大重试次数
	InitialDelay = 2 * time.Second // 初始重试间隔
)

// updatePostVoteCount 更新帖子的投票数,同时更新 Redis 的帖子热度
// caseNum 用于区分不同的投票情况
// 1. 之前没投过票，现在要投赞成票
// 2. 之前投过反对票，现在要改为赞成票
// 3. 之前投过赞成票，现在要取消
// 4. 之前投过反对票，现在要取消
// 5. 之前没投过票，现在要投反对票
// 6. 之前投过赞成票，现在要改为反对票
// 通过 caseNum 调用不同的 dao 方法
// updatePostVoteCount 更新帖子投票数并更新 Redis 热度
func updatePostVoteCount(ctx context.Context, v *VoteService, id int64, voteFor, caseNum int) {
	// 创建带超时的 context，避免无限重试
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 定义 caseNum 和对应 DAO 方法的映射
	caseFuncMap := map[int]func(context.Context, int64, int) error{
		1: v.VoteDaoInterface.ContentVoteCase1,
		2: v.VoteDaoInterface.ContentVoteCase2,
		3: v.VoteDaoInterface.ContentVoteCase3,
		4: v.VoteDaoInterface.ContentVoteCase4,
		5: v.VoteDaoInterface.ContentVoteCase5,
		6: v.VoteDaoInterface.ContentVoteCase6,
	}

	// 检查是否有对应的处理函数
	voteFunc, ok := caseFuncMap[caseNum]
	if !ok {
		zap.L().Error("Invalid case number", zap.Int("caseNum", caseNum))
		return
	}

	// 执行带重试机制的投票逻辑
	var attempt int
	var err error
	delay := InitialDelay

	for attempt = 1; attempt <= MaxRetries; attempt++ {
		err = voteFunc(ctx, id, voteFor) // 执行对应的 DAO 方法
		if err == nil {
			break
		}
		log.Printf("Attempt %d failed: %v. Retrying...", attempt, err)
		time.Sleep(delay)
		delay *= 2 // 指数退避，增加重试间隔
	}

	if err != nil {
		zap.L().Error("Failed to update post vote count", zap.Error(err))
		return
	}

	// 成功更新投票后，查询新的票数和创建时间
	if voteFor == dao.VotePost {
		up, down, err := v.VoteDaoInterface.GetContentVoteCount(ctx, id)
		if err != nil {
			log.Printf("Failed to get vote count: %v", err)
			return
		}

		createTime, err := v.VoteDaoInterface.GetPostCreateTime(ctx, id)
		if err != nil {
			zap.L().Error("Failed to get post create time", zap.Error(err))
			return
		}

		// 更新 Redis 中的帖子热度
		if err := v.VoteCacheInterface.UpdatePostHot(ctx, id, int(up), int(down), createTime); err != nil {
			zap.L().Error("Failed to update post hot score", zap.Error(err))
		}
	}
}
