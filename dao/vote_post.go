package dao

import (
	"GinTalk/DTO"
	"GinTalk/dao/MySQL"
	"GinTalk/model"
	"context"
	"fmt"
	"time"
)

func AddPostVote(ctx context.Context, postID int64, userID int64) error {
	vote := model.VotePost{
		PostID: postID,
		UserID: userID,
		Vote:   1,
	}
	sqlStr := `
		INSERT INTO vote_post (post_id, user_id, vote)
		VALUES (?, ?, ?)	
`
	return MySQL.GetDB().WithContext(ctx).Exec(sqlStr, vote.PostID, vote.UserID, vote.Vote).Error
}

func RevokePostVote(ctx context.Context, postID int64, userID int64) error {
	sqlStr := `
		DELETE FROM vote_post
		WHERE post_id = ? AND user_id = ?`
	return MySQL.GetDB().WithContext(ctx).Exec(sqlStr, postID, userID).Error
}

func AddContentVoteUp(ctx context.Context, postID int64) error {
	sqlStr := `
		UPDATE content_votes
		SET vote = vote + 1
		WHERE post_id = ? AND delete_time = 0`
	return MySQL.GetDB().WithContext(ctx).Exec(sqlStr, postID).Error
}

func SubContentVoteUp(ctx context.Context, postID int64) error {
	sqlStr := `
		UPDATE content_votes
		SET vote = vote - 1
		WHERE post_id = ? AND delete_time = 0`
	return MySQL.GetDB().WithContext(ctx).Exec(sqlStr, postID).Error
}

func GetUserVoteList(ctx context.Context, userID int64, pageNum int, pageSize int) ([]int64, error) {
	var voteRecord []int64
	sqlStr := `
		SELECT post_id
		FROM vote_post
		WHERE user_id = ?
		LIMIT ? OFFSET ?`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, userID, pageSize, (pageNum-1)*pageSize).Scan(&voteRecord).Error
	return voteRecord, err
}

func GetPostVoteCount(ctx context.Context, postID int64) (*DTO.PostVoteCounts, error) {
	var voteCount DTO.PostVoteCounts
	sqlStr := `
		SELECT post_id, vote
		FROM content_votes
		WHERE post_id = ? AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID).Scan(&voteCount).Error
	return &voteCount, err
}

func GetBatchPostVoteCount(ctx context.Context, postIDs []int64) ([]DTO.PostVoteCounts, error) {
	var voteCount []DTO.PostVoteCounts
	sqlStr := `
		SELECT post_id, vote
		FROM content_votes
		WHERE post_id IN (?) AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postIDs).Scan(&voteCount).Error
	return voteCount, err
}


// CheckUserVoted 检查用户是否对一组帖子投票。
// 它查询数据库中与给定帖子ID和用户ID相关的投票记录，
// 并返回一个包含帖子ID和投票信息的 UserVotePostRelationsDTO 切片。
//
// 参数:
//   - ctx: 用于管理请求范围内的值、取消和截止日期的上下文。
//   - postIDs: 要检查投票的帖子ID切片。
//   - userID: 要检查投票的用户ID。
//
// 返回:
//   - 一个包含帖子ID和投票信息的 UserVotePostRelationsDTO 切片。
//   - 如果查询失败或执行过程中出现任何其他问题，则返回错误。
func CheckUserVoted(ctx context.Context, postIDs []int64, userID int64) ([]DTO.UserVotePostRelationsDTO, error) {
	var votes []DTO.UserVotePostRelationsDTO
	sqlStr := `
		SELECT post_id, vote
		FROM vote_post
		WHERE post_id IN (?) AND user_id = ? AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postIDs, userID).Scan(&votes).Error
	return votes, err
}

// GetPostVoteDetail 获取特定帖子的投票详情。
//
// 该函数通过分页方式查询指定帖子的投票信息，包括投票用户的ID、用户名和投票值。
// 查询结果会排除已删除的投票记录和用户。
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期。
//   - postID: 要查询的帖子ID。
//   - pageNum: 页码，用于分页查询，从1开始。
//   - pageSize: 每页的记录数，用于分页查询。
//
// 返回:
//   - []DTO.UserVotePostDetailDTO: 包含投票详情的切片，每个元素代表一条投票记录。
//   - error: 如果查询过程中发生错误，将返回相应的错误信息；否则返回nil。
func GetPostVoteDetail(ctx context.Context, postID int64, pageNum int, pageSize int) ([]DTO.UserVotePostDetailDTO, error) {
	var votes []DTO.UserVotePostDetailDTO
	sqlStr := `
		SELECT user.user_id, post_id, vote, username
		FROM vote_post
		INNER JOIN user ON vote_post.user_id = user.user_id
		WHERE post_id = ? AND vote_post.delete_time = 0 AND user.delete_time = 0
		LIMIT ? OFFSET ?`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID, pageSize, (pageNum-1)*pageSize).Scan(&votes).Error
	return votes, err
}

// GetPostCreateTime 从数据库中检索帖子的创建时间。
// 它接受一个上下文和一个帖子 ID 作为参数，并返回创建时间和查询过程中发生的错误（如果有）。
//
// 参数:
//   - ctx: 用于管理请求范围内的值、取消和截止日期的上下文。
//   - postID: 要检索创建时间的帖子的 ID。
//
// 返回:
//   - time.Time: 指定帖子的创建时间。
//   - error: 如果查询过程中发生错误，则返回错误对象，否则返回 nil。
func GetPostCreateTime(ctx context.Context, postID int64) (time.Time, error) {
	var createTime time.Time
	sqlStr := `
		SELECT create_time
		FROM post
		WHERE post_id = ?`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID).Scan(&createTime).Error
	return createTime, err
}

// CheckUserVotedPost 检查用户是否对特定帖子投票。
// 如果用户已投票，则返回 true，否则返回 false。
//
// 参数:
// - ctx: 用于管理请求范围内的值、取消和截止日期的上下文。
// - postID: 要检查的帖子的 ID。
// - userID: 要检查的用户的 ID。
//
// 返回:
// - bool: 如果用户已对帖子投票，则为 true，否则为 false。
// - error: 查询执行过程中遇到的任何错误。
func CheckUserVotedPost(ctx context.Context, postID int64, userID int64) (bool, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*)
		FROM vote_post
		WHERE post_id = ? AND user_id = ? AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID, userID).Scan(&count).Error
	return count > 0, err
}

// AddPostVoteWithTx 在事务中添加或删除帖子投票。
// 它根据投票值在 vote_post 表中插入或删除记录，并相应地更新 content_votes 表中的投票计数。
//
// 参数:
//   - ctx: 用于管理请求范围内的值、取消和截止日期的上下文。
//   - postID: 要投票的帖子的 ID。
//   - userID: 投票用户的 ID。
//   - vote: 投票值，正值表示赞成票，非正值表示反对票。
//
// 返回:
//   - error: 如果事务失败则返回错误，否则返回 nil。
func AddPostVoteWithTx(ctx context.Context, postID int64, userID int64, vote int) error {
	tx := MySQL.GetDB().WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var sqlStr string
	if vote > 0 {
		sqlStr = `
		INSERT INTO vote_post (post_id, user_id)
		VALUES (?, ?)	
`
	} else {
		sqlStr = `
		DELETE FROM vote_post
		WHERE post_id = ? AND user_id = ?`
	}
	result := tx.Exec(sqlStr, postID, userID)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("vote failed")
	}
	if vote > 0 {
		sqlStr = `
		UPDATE content_votes
		SET vote = vote + 1
		WHERE post_id = ? AND delete_time = 0`
	} else {
		sqlStr = `
		UPDATE content_votes
		SET vote = vote - 1
		WHERE post_id = ? AND delete_time = 0`
	}
	if err := tx.Exec(sqlStr, postID).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
