package dao

import (
	"GinTalk/dao/MySQL"
	"GinTalk/model"
	"context"
	"time"
)

// GetTopComments retrieves the top-level comments for a post.
// It fetches comments that are not deleted and have a status of 1,
// ordered by creation time in descending order. Pagination is supported
// through pageSize and pageNum parameters.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation, and deadlines.
//   - postID: The ID of the post for which to retrieve comments.
//   - pageSize: The number of comments to retrieve per page.
//   - pageNum: The page number to retrieve.
//
// Returns:
//   - A slice of model.Comment containing the retrieved comments.
//   - An error if the operation fails.
func GetTopComments(ctx context.Context, postID int64, pageSize, pageNum int) ([]model.Comment, error) {
	var comment []model.Comment
	sqlStr := `
		SELECT * 
		FROM comment
		INNER JOIN comment_relation ON comment.comment_id = comment_relation.comment_id
		WHERE comment.post_id = ? AND comment.status = 1 AND comment.delete_time = 0 AND comment_relation.delete_time = 0 AND comment_relation.parent_id = 0
		ORDER BY comment.create_time DESC
		LIMIT ? OFFSET ?
		`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID, pageSize, (pageNum-1)*pageSize).Scan(&comment).Error

	return comment, err
}

// GetSubComments 获取某一个顶层评论的子评论
func GetSubComments(ctx context.Context, postID, parentID int64, pageSize, pageNum int) ([]model.Comment, error) {
	var comments []model.Comment
	sqlStr := `
		SELECT * 
		FROM comment
		INNER JOIN comment_relation ON comment.comment_id = comment_relation.comment_id
		WHERE comment.post_id = ? AND parent_id = ? AND status = 1 AND comment.delete_time = 0 AND comment_relation.delete_time = 0
		ORDER BY comment.create_time DESC
		LIMIT ? OFFSET ?`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID, parentID, pageSize, (pageNum-1)*pageSize).Scan(&comments).Error

	return comments, err
}

// GetCommentByID 根据评论 ID 获取评论
func GetCommentByID(ctx context.Context, commentID int64) (*model.Comment, error) {
	var comment model.Comment
	sqlStr := `
		SELECT * FROM comment
		WHERE comment_id = ? AND status = 1 AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, commentID).Scan(&comment).Error
	return &comment, err
}

func GetCommentRelationByID(ctx context.Context, commentID int64) (*model.CommentRelation, error) {
	var relation model.CommentRelation
	sqlStr := `
		SELECT * FROM comment_relation
		WHERE comment_id = ? AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, commentID).Scan(&relation).Error
	return &relation, err
}

func GetCommentParentUserID(ctx context.Context, commentID int64) (int64, error) {
	var userID int64
	sqlStr := `
		SELECT author_id FROM comment
		INNER JOIN comment_relation ON comment.comment_id = comment_relation.parent_id
		WHERE comment_relation.comment_id = ?`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, commentID).Scan(&userID).Error
	return userID, err
}

// CreateComment 创建评论
func CreateComment(ctx context.Context, comment *model.Comment, replyID int64, parentID int64) error {
	tx := MySQL.GetDB().Begin().WithContext(ctx)
	sqlStrCreateComment := `
		INSERT INTO comment (comment_id, content, post_id, author_id, author_name)
			VALUES (?, ?, ?, ?, ?)`
	err := tx.Exec(sqlStrCreateComment, comment.CommentID, comment.Content, comment.PostID, comment.AuthorID, comment.AuthorName).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	sqlStrCreateRelation := `
		INSERT INTO comment_relation (post_id, comment_id, parent_id, reply_id) 
			VALUES (?, ?, ?, ?)`
	err = tx.Exec(sqlStrCreateRelation, comment.PostID, comment.CommentID, parentID, replyID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// UpdateComment 更新评论
func UpdateComment(ctx context.Context, commentID int64, content string) error {
	sqlStr := `
		UPDATE comment
		SET content = ?
		WHERE comment_id = ?`
	return MySQL.GetDB().WithContext(ctx).Exec(sqlStr, content, commentID).Error
}

// DeleteComment 删除评论
func DeleteComment(ctx context.Context, commentID int64) error {
	sqlStrDeleteComment := `
		UPDATE comment
		SET delete_time = ?
		WHERE comment_id = ?`
	sqlStrDeleteRelation := `
		UPDATE comment_relation
		SET delete_time = ?
		WHERE comment_id = ? OR parent_id = ? OR reply_id = ?`
	tx := MySQL.GetDB().Begin().WithContext(ctx)
	err := tx.Exec(sqlStrDeleteComment, time.Now().Unix(), commentID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Exec(sqlStrDeleteRelation, time.Now().Unix(), commentID, commentID, commentID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// GetCommentCount 获取评论数量
func GetCommentCount(ctx context.Context, postID int64) (int64, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*) FROM comment
		WHERE post_id = ? AND status = 1 AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID).Scan(&count).Error
	return count, err
}

// GetTopCommentCount 获取顶级评论数量
func GetTopCommentCount(ctx context.Context, postID int64) (int64, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*) FROM comment
		INNER JOIN comment_relation ON comment.comment_id = comment_relation.comment_id
		WHERE comment.post_id = ? AND status = 1 AND comment.delete_time = 0 AND parent_id = 0 AND comment_relation.delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, postID).Scan(&count).Error
	return count, err
}

// GetSubCommentCount 获取子评论数量
func GetSubCommentCount(ctx context.Context, parentID int64) (int64, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*) FROM comment
		INNER JOIN comment_relation ON comment.comment_id = comment_relation.parent_id
		WHERE parent_id = ? AND status = 1 AND comment.delete_time = 0 AND comment_relation.delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, parentID).Scan(&count).Error
	return count, err
}

// GetCommentCountByUserID 获取用户评论数量
func GetCommentCountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*) FROM comment
		WHERE author_id = ? AND status = 1 AND delete_time = 0`
	err := MySQL.GetDB().WithContext(ctx).Raw(sqlStr, userID).Scan(&count).Error
	return count, err
}
