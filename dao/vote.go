package dao

import (
	"GinTalk/model"
	"context"
	"gorm.io/gorm"
)

/*
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

type VoteDaoInterface interface {
	// VoteCase1 之前没投过票，现在要投赞成票
	VoteCase1(ctx context.Context, postID int64, userID int64) error
	ContentVoteCase1(ctx context.Context, postID int64) error
	// VoteCase2 之前投过反对票，现在要改为赞成票
	VoteCase2(ctx context.Context, postID int64, userID int64) error
	ContentVoteCase2(ctx context.Context, postID int64) error
	// VoteCase3 之前投过赞成票，现在要取消
	VoteCase3(ctx context.Context, postID int64, userID int64) error
	ContentVoteCase3(ctx context.Context, postID int64) error
	// VoteCase4 之前投过反对票，现在要取消
	VoteCase4(ctx context.Context, postID int64, userID int64) error
	ContentVoteCase4(ctx context.Context, postID int64) error
	// VoteCase5 之前没投过票，现在要投反对票
	VoteCase5(ctx context.Context, postID int64, userID int64) error
	ContentVoteCase5(ctx context.Context, postID int64) error
	// VoteCase6 之前投过赞成票，现在要改为反对票
	VoteCase6(ctx context.Context, postID int64, userID int64) error
	ContentVoteCase6(ctx context.Context, postID int64) error

	// GetVoteRecord 获取用户对某个帖子的投票记录
	GetVoteRecord(ctx context.Context, postID int64, userID int64) (int, error)

	// RevokeVote 取消投票
	//RevokeVote(ctx context.Context, postID int64, userID int64) error

	// AddContentVoteUp 帖子赞成票数加一
	AddContentVoteUp(ctx context.Context, postID int64) error
	// SubContentVoteUp 帖子赞成票数减一
	SubContentVoteUp(ctx context.Context, postID int64) error
	// AddContentVoteDown 帖子反对票数加一
	AddContentVoteDown(ctx context.Context, postID int64) error
	// SubContentVoteDown 帖子反对票数减一
	SubContentVoteDown(ctx context.Context, postID int64) error

	// GetUserVoteList 获取用户投票过的帖子
	GetUserVoteList(ctx context.Context, userID int64, pageNum int, pageSize int) ([]int64, error)

	// GetContentVoteCount 获取帖子的赞成票数和反对票数
	GetContentVoteCount(ctx context.Context, postID int64) (int64, int64, error)

	// CheckUserVoted 检查用户是否投过票
	CheckUserVoted(ctx context.Context, postID []int64, userID int64) ([]model.Vote, error)
}

type VoteDao struct {
	*gorm.DB
}

func NewVoteDao(db *gorm.DB) VoteDaoInterface {
	return &VoteDao{DB: db}
}

func (vd *VoteDao) VoteCase1(ctx context.Context, postID int64, userID int64) error {
	sqlStr := `INSERT INTO vote (post_id, user_id, vote) VALUES (?, ?, 1)`
	return vd.WithContext(ctx).Exec(sqlStr, postID, userID).Error
}

func (vd *VoteDao) ContentVoteCase1(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET up = up + 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) VoteCase2(ctx context.Context, postID int64, userID int64) error {
	sqlStr := `UPDATE vote SET vote = 1 WHERE post_id = ? AND user_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID, userID).Error
}

func (vd *VoteDao) ContentVoteCase2(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET up = up + 1, down = down - 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) VoteCase3(ctx context.Context, postID int64, userID int64) error {
	sqlStr := `UPDATE vote SET delete_time = CURRENT_TIME WHERE post_id = ? AND user_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID, userID).Error
}

func (vd *VoteDao) ContentVoteCase3(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET up = up - 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) VoteCase4(ctx context.Context, postID int64, userID int64) error {
	sqlStr := `UPDATE vote SET delete_time = CURRENT_TIME WHERE post_id = ? AND user_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID, userID).Error
}

func (vd *VoteDao) ContentVoteCase4(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET down = down - 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) VoteCase5(ctx context.Context, postID int64, userID int64) error {
	sqlStr := `INSERT INTO vote (post_id, user_id, vote) VALUES (?, ?, -1)`
	return vd.WithContext(ctx).Exec(sqlStr, postID, userID).Error
}

func (vd *VoteDao) ContentVoteCase5(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET down = down + 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) VoteCase6(ctx context.Context, postID int64, userID int64) error {
	sqlStr := `UPDATE vote SET vote = -1 WHERE post_id = ? AND user_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID, userID).Error
}

func (vd *VoteDao) ContentVoteCase6(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET up = up - 1, down = down + 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) GetVoteRecord(ctx context.Context, postID int64, userID int64) (int, error) {
	sqlStr := `SELECT vote FROM vote WHERE post_id = ? AND user_id = ? AND delete_time = 0`
	var vote int
	err := vd.WithContext(ctx).Raw(sqlStr, postID, userID).Scan(&vote).Error
	return vote, err
}

func (vd *VoteDao) AddContentVoteUp(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET up = up + 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) SubContentVoteUp(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET up = up - 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) AddContentVoteDown(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET down = down + 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) SubContentVoteDown(ctx context.Context, postID int64) error {
	sqlStr := `UPDATE content_votes SET down = down - 1 WHERE post_id = ? AND delete_time = 0`
	return vd.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (vd *VoteDao) GetUserVoteList(ctx context.Context, userID int64, pageNum int, pageSize int) ([]int64, error) {
	sqlStr := `SELECT post_id FROM vote WHERE user_id = ? AND delete_time = 0 LIMIT ? OFFSET ?`
	var postIDs []int64
	err := vd.WithContext(ctx).Raw(sqlStr, userID, pageSize, (pageNum-1)*pageSize).Scan(&postIDs).Error
	return postIDs, err
}

func (vd *VoteDao) GetContentVoteCount(ctx context.Context, postID int64) (int64, int64, error) {
	type votes struct {
		Up   int64 `db:"up"`
		Down int64 `db:"down"`
	}
	var vote votes
	sqlStr := `SELECT up, down FROM content_votes WHERE post_id = ? AND delete_time = 0`
	err := vd.WithContext(ctx).Raw(sqlStr, postID).Scan(&vote).Error
	return vote.Up, vote.Down, err
}

func (vd *VoteDao) CheckUserVoted(ctx context.Context, postID []int64, userID int64) ([]model.Vote, error) {
	var votes []model.Vote
	sqlStr := `SELECT post_id, user_id, vote FROM vote WHERE post_id IN (?) AND user_id = ? AND delete_time IS NULL`
	err := vd.WithContext(ctx).Raw(sqlStr, postID, userID).Scan(&votes).Error
	return votes, err
}
