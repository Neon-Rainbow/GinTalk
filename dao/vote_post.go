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

func CheckUserVoted(ctx context.Context, postIDs []int64, userID int64) ([]DTO.UserVotePostRelationsDTO, error) {
	var votes []DTO.UserVotePostRelationsDTO
	sqlStr := `
		SELECT post_id, vote
		FROM vote_post
		WHERE post_id IN (?) AND user_id = ? AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postIDs, userID).Scan(&votes).Error
	return votes, err
}

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

func GetPostCreateTime(ctx context.Context, postID int64) (time.Time, error) {
	var createTime time.Time
	sqlStr := `
		SELECT create_time
		FROM post
		WHERE post_id = ?`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID).Scan(&createTime).Error
	return createTime, err
}

func CheckUserVotedPost(ctx context.Context, postID int64, userID int64) (bool, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*)
		FROM vote_post
		WHERE post_id = ? AND user_id = ? AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID, userID).Scan(&count).Error
	return count > 0, err
}

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
