package service

import (
	"GinTalk/dao"
	"GinTalk/dao/MySQL"
	"GinTalk/model"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"context"
	"fmt"
	"time"
)

var _ VoteServiceInterface = (*VoteService)(nil)

type VoteServiceInterface interface {
	Vote(ctx context.Context, postID int64, userID int64, voteType int) *apiError.ApiError
	RevokeVote(ctx context.Context, postID int64, userID int64) *apiError.ApiError
	MyVoteList(ctx context.Context, userID int64, pageNum int, pageSize int) ([]int64, *apiError.ApiError)
	GetVoteCount(ctx context.Context, postID int64) (int64, int64, *apiError.ApiError)
	CheckUserVoted(ctx context.Context, postID []int64, userID int64) ([]model.Vote, *apiError.ApiError)
}

type VoteService struct {
	voteDao dao.IVoteDo
}

func (v *VoteService) Vote(ctx context.Context, postID int64, userID int64, voteType int) *apiError.ApiError {
	// 查询先前的投票记录
	voteRecord, err := v.voteDao.WithContext(ctx).
		Where(dao.Vote.PostID.Eq(postID), dao.Vote.UserID.Eq(userID)).
		First()
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票记录失败",
		}
	}
	// 用于记录帖子的投票数的变化量
	voteChange := 0

	if voteRecord == nil {
		voteChange = voteType
		// 创建新的投票记录
		err = v.voteDao.WithContext(ctx).Create(&model.Vote{
			PostID: postID,
			UserID: userID,
			Vote:   int32(voteType),
		})
	} else {
		// 更新投票记录
		voteChange = voteType - int(voteRecord.Vote)
		voteRecord.Vote = int32(voteType)
		err = v.voteDao.WithContext(ctx).Save(voteRecord)
	}
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "保存投票记录失败",
		}
	}
	// 异步更新帖子的投票数
	go func() {
		if voteChange == 0 {
			return
		}

		var maxRetries = 3          // 最大重试次数
		var delay = 2 * time.Second // 重试间隔时间

		for attempt := 1; attempt <= maxRetries; attempt++ {
			switch voteChange {
			case 2:
				// 开启事务
				tx := MySQL.GetDB().Begin()
				// 将帖子的up + 1
				err := tx.WithContext(ctx).Raw("UPDATE post SET up = up + 1 WHERE post_id = ?", postID).Error
				if err != nil {
					tx.Rollback()
					time.Sleep(delay)
					continue
				}
				// 将帖子的down - 1
				err = tx.WithContext(ctx).Raw("UPDATE post SET down = down - 1 WHERE post_id = ?", postID).Error
				if err != nil {
					tx.Rollback()
					time.Sleep(delay)
					continue
				}
				// 提交事务
				tx.Commit()
				return
			case 1:
				err := MySQL.GetDB().Raw("UPDATE post SET up = up + 1 WHERE post_id = ?", postID).Error
				if err != nil {
					time.Sleep(delay)
					continue
				}
				return
			case -1:
				err := MySQL.GetDB().Raw("UPDATE post SET down = down + 1 WHERE post_id = ?", postID).Error
				if err != nil {
					time.Sleep(delay)
					continue
				}
				return
			case -2:
				// 开启事务
				tx := MySQL.GetDB().Begin()
				// 将帖子的up - 1
				err := tx.WithContext(ctx).Raw("UPDATE post SET up = up - 1 WHERE post_id = ?", postID).Error
				if err != nil {
					tx.Rollback()
					time.Sleep(delay)
					continue
				}
				// 将帖子的down + 1
				err = tx.WithContext(ctx).Raw("UPDATE post SET down = down + 1 WHERE post_id = ?", postID).Error
				if err != nil {
					tx.Rollback()
					time.Sleep(delay)
					continue
				}
				return
			}
		}
	}()
	return nil
}

func (v *VoteService) RevokeVote(ctx context.Context, postID int64, userID int64) *apiError.ApiError {
	// 查看原先的投票记录
	voteRecord, err := v.voteDao.WithContext(ctx).
		Where(dao.Vote.PostID.Eq(postID), dao.Vote.UserID.Eq(userID)).
		First()
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票记录失败",
		}
	}
	if voteRecord == nil {
		return nil
	}
	// 删除投票记录
	_, err = v.voteDao.WithContext(ctx).Delete(voteRecord)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "删除投票记录失败",
		}
	}
	// 异步更新帖子的投票数
	go func() {
		var maxRetries = 3          // 最大重试次数
		var delay = 2 * time.Second // 重试间隔时间

		for attempt := 1; attempt <= maxRetries; attempt++ {
			switch voteRecord.Vote {
			case 1:
				err := MySQL.GetDB().Raw("UPDATE post SET up = up - 1 WHERE post_id = ?", postID).Error
				if err != nil {
					time.Sleep(delay)
					continue
				}
				return
			case -1:
				err := MySQL.GetDB().Raw("UPDATE post SET down = down - 1 WHERE post_id = ?", postID).Error
				if err != nil {
					time.Sleep(delay)
					continue
				}
				return
			}
		}
	}()
	return nil
}

func (v *VoteService) MyVoteList(ctx context.Context, userID int64, pageNum int, pageSize int) ([]int64, *apiError.ApiError) {
	voteRecord := make([]model.Vote, pageSize)
	err := v.voteDao.WithContext(ctx).
		Where(dao.Vote.UserID.Eq(userID)).
		Order(dao.Vote.CreateTime.Desc()).
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Scan(&voteRecord)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票记录失败",
		}
	}
	postIDList := make([]int64, len(voteRecord))
	for i, v := range voteRecord {
		postIDList[i] = v.PostID
	}
	return postIDList, nil
}

func (v *VoteService) GetVoteCount(ctx context.Context, postID int64) (int64, int64, *apiError.ApiError) {
	// 查询帖子的投票数
	type votes struct {
		Up   int64 `db:"up"`
		Down int64 `db:"down"`
	}
	var vote votes
	err := MySQL.GetDB().WithContext(ctx).Model(&model.Post{}).Select("up", "down").Where("post_id = ?", postID).Scan(&vote).Error
	if err != nil {
		return 0, 0, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "查询投票数失败",
		}
	}
	return vote.Up, vote.Down, nil
}

// CheckUserVoted 批量查询用户是否投票过
func (v *VoteService) CheckUserVoted(ctx context.Context, postID []int64, userID int64) ([]model.Vote, *apiError.ApiError) {
	var votes []model.Vote

	// 构建原生 SQL 查询
	sqlStr := `SELECT post_id, user_id, vote
               FROM vote 
               WHERE post_id IN (?) AND delete_time IS NULL AND user_id = ?`

	// 使用 Raw() 执行查询
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID, userID).Scan(&votes).Error
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  fmt.Sprintf("批量查询投票记录失败: %v", err),
		}
	}
	return votes, nil
}

func NewVoteService(voteDao dao.IVoteDo) VoteServiceInterface {
	return &VoteService{
		voteDao: voteDao,
	}
}
