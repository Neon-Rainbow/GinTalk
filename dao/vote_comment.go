package dao

import (
	"gorm.io/gorm"
)

var _ CommentVoteInterface = (*CommentVoteImpl)(nil)

type CommentVoteInterface interface {
	// VoteComment 为评论投票
	VoteComment(userID, commentID int64) error

	// RemoveVoteComment 取消评论投票
	RemoveVoteComment(userID, commentID int64) error

	// GetVoteComment 获取评论投票数
	GetVoteComment(userID, commentID int64) (int, error)

	// GetCommentVoteStatus 获取用户对该评论的投票状态
	GetCommentVoteStatus(userID, commentID int64) (int, error)

	// GetCommentVoteStatusList 用户批量查询是否投票
	GetCommentVoteStatusList(userID int64, commentIDs []int64) (map[int64]int, error)

	// IncrCommentVoteCount 增加评论投票数
	IncrCommentVoteCount(commentID int64) error

	// DecrCommentVoteCount 减少评论投票数
	DecrCommentVoteCount(commentID int64) error
}

type CommentVoteImpl struct {
	*gorm.DB
}

func NewCommentVoteImpl(db *gorm.DB) CommentVoteInterface {
	return &CommentVoteImpl{DB: db}
}

func (v *CommentVoteImpl) VoteComment(userID, commentID int64) error {
	sqlStr := `
	INSERT INTO vote_comment (user_id, comment_id, vote) 
	VALUES (?, ?, 1)`
	return v.Exec(sqlStr, userID, commentID).Error
}

func (v *CommentVoteImpl) RemoveVoteComment(userID, commentID int64) error {
	sqlStr := `
	DELETE FROM vote_comment
	WHERE user_id = ? AND comment_id = ?`
	return v.Exec(sqlStr, userID, commentID).Error
}

func (v *CommentVoteImpl) GetVoteComment(userID, commentID int64) (int, error) {
	var count int
	sqlStr := `
	SELECT COUNT(*) 
	FROM vote_comment
	WHERE comment_id = ? AND vote = 1`
	err := v.Raw(sqlStr, commentID).Scan(&count).Error
	return count, err
}

func (v *CommentVoteImpl) GetCommentVoteStatus(userID, commentID int64) (int, error) {
	var vote int
	sqlStr := `
	SELECT vote
	FROM vote_comment
	WHERE user_id = ? AND comment_id = ?`
	err := v.Raw(sqlStr, userID, commentID).Scan(&vote).Error
	return vote, err
}

func (v *CommentVoteImpl) GetCommentVoteStatusList(userID int64, commentIDs []int64) (map[int64]int, error) {
	voteMap := make(map[int64]int)
	sqlStr := `
	SELECT comment_id, vote
	FROM vote_comment
	WHERE user_id = ? AND comment_id IN (?)`
	rows, err := v.Raw(sqlStr, userID, commentIDs).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var commentID, vote int64
		rows.Scan(&commentID, &vote)
		voteMap[commentID] = int(vote)
	}
	return voteMap, nil
}

func (v *CommentVoteImpl) IncrCommentVoteCount(commentID int64) error {
	sqlStr := `
	UPDATE comment_votes
	SET up = up + 1
	WHERE comment_id = ? AND delete_time = 0`
	return v.Exec(sqlStr, commentID).Error
}

func (v *CommentVoteImpl) DecrCommentVoteCount(commentID int64) error {
	sqlStr := `
	UPDATE comment_votes
	SET up = up - 1
	WHERE comment_id = ? AND delete_time = 0`
	return v.Exec(sqlStr, commentID).Error
}
