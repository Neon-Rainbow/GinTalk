package dao

import (
	"GinTalk/dao/MySQL"
)

func VoteComment(userID, commentID int64) error {
	sqlStr := `
	INSERT INTO vote_comment (user_id, comment_id, vote) 
	VALUES (?, ?, 1)`
	return MySQL.GetDB().Exec(sqlStr, userID, commentID).Error
}

func RemoveVoteComment(userID, commentID int64) error {
	sqlStr := `
	DELETE FROM vote_comment
	WHERE user_id = ? AND comment_id = ?`
	return MySQL.GetDB().Exec(sqlStr, userID, commentID).Error
}

func GetVoteComment(userID, commentID int64) (int, error) {
	var count int
	sqlStr := `
	SELECT COUNT(*) 
	FROM vote_comment
	WHERE comment_id = ? AND vote = 1`
	err := MySQL.GetDB().Raw(sqlStr, commentID).Scan(&count).Error
	return count, err
}

func GetCommentVoteStatus(userID, commentID int64) (int, error) {
	var vote int
	sqlStr := `
	SELECT vote
	FROM vote_comment
	WHERE user_id = ? AND comment_id = ?`
	err := MySQL.GetDB().Raw(sqlStr, userID, commentID).Scan(&vote).Error
	return vote, err
}

func GetCommentVoteStatusList(userID int64, commentIDs []int64) (map[int64]int, error) {
	voteMap := make(map[int64]int)
	sqlStr := `
	SELECT comment_id, vote
	FROM vote_comment
	WHERE user_id = ? AND comment_id IN (?)`
	rows, err := MySQL.GetDB().Raw(sqlStr, userID, commentIDs).Rows()
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

func IncrCommentVoteCount(commentID int64) error {
	sqlStr := `
	UPDATE comment_votes
	SET up = up + 1
	WHERE comment_id = ? AND delete_time = 0`
	return MySQL.GetDB().Exec(sqlStr, commentID).Error
}

func DecrCommentVoteCount(commentID int64) error {
	sqlStr := `
	UPDATE comment_votes
	SET up = up - 1
	WHERE comment_id = ? AND delete_time = 0`
	return MySQL.GetDB().Exec(sqlStr, commentID).Error
}
