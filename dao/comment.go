package dao

import (
	"GinTalk/model"
	"context"
	"gorm.io/gorm"
	"time"
)

type CommentDaoInterface interface {
	// GetTopComments 获取顶层评论
	GetTopComments(ctx context.Context, postID int64, pageSize, pageNum int) ([]*model.Comment, error)
	// GetSubComments 获取某一个顶层评论的子评论
	GetSubComments(ctx context.Context, postID, parentID int64, pageSize, pageNum int) ([]*model.Comment, error)
	// GetCommentByID 根据评论 ID 获取评论
	GetCommentByID(ctx context.Context, commentID int64) (*model.Comment, error)
	// CreateComment 创建评论
	CreateComment(ctx context.Context, comment *model.Comment, replyID int64, parentID int64) error
	// UpdateComment 更新评论
	UpdateComment(ctx context.Context, commentID int64, content string) error
	// DeleteComment 删除评论
	DeleteComment(ctx context.Context, commentID int64) error
	// GetCommentCount 获取评论数量
	GetCommentCount(ctx context.Context, postID int64) (int64, error)
	// GetTopCommentCount 获取某条内容下的顶级评论数量
	GetTopCommentCount(ctx context.Context, postID int64) (int64, error)
	// GetSubCommentCount 获取某条一级评论下的二级评论数量
	GetSubCommentCount(ctx context.Context, parentID int64) (int64, error)
	// GetCommentCountByUserID 获取用户评论数量
	GetCommentCountByUserID(ctx context.Context, userID int64) (int64, error)
}

type CommentDao struct {
	*gorm.DB
}

func NewCommentDao(db *gorm.DB) CommentDaoInterface {
	return &CommentDao{DB: db}
}

// GetTopComments 获取顶层评论
func (cd *CommentDao) GetTopComments(ctx context.Context, postID int64, pageSize, pageNum int) ([]*model.Comment, error) {
	var comment []*model.Comment
	sqlStr := `
		SELECT * 
		FROM comment
		INNER JOIN comment_relation ON comment.comment_id = comment_relation.comment_id
		WHERE comment.post_id = ? AND comment.status = 1 AND comment.delete_time = 0 AND comment_relation.delete_time = 0 AND comment_relation.parent_id = 0
		ORDER BY comment.create_time DESC`
	err := cd.WithContext(ctx).Raw(sqlStr, postID).Scan(&comment).Error
	return comment, err
}

// GetSubComments 获取某一个顶层评论的子评论
func (cd *CommentDao) GetSubComments(ctx context.Context, postID, parentID int64, pageSize, pageNum int) ([]*model.Comment, error) {
	var comments []*model.Comment
	sqlStr := `
		SELECT * 
		FROM comment
		INNER JOIN comment_relation ON comment.comment_id = comment_relation.comment_id
		WHERE comment.post_id = ? AND parent_id = ? AND status = 1 AND comment.delete_time = 0 AND comment_relation.delete_time = 0
		ORDER BY comment.create_time DESC
		LIMIT ? OFFSET ?`
	err := cd.WithContext(ctx).Raw(sqlStr, postID, parentID, pageSize, (pageNum-1)*pageSize).Scan(&comments).Error
	return comments, err
}

// GetCommentByID 根据评论 ID 获取评论
func (cd *CommentDao) GetCommentByID(ctx context.Context, commentID int64) (*model.Comment, error) {
	var comment model.Comment
	sqlStr := `
		SELECT * FROM comment
		WHERE comment_id = ? AND status = 1 AND delete_time = 0`
	err := cd.WithContext(ctx).Raw(sqlStr, commentID).Scan(&comment).Error
	return &comment, err
}

// CreateComment 创建评论
func (cd *CommentDao) CreateComment(ctx context.Context, comment *model.Comment, replyID int64, parentID int64) error {
	tx := cd.Begin().WithContext(ctx)
	sqlStrCreateComment := `
		INSERT INTO comment (comment_id, content, post_id, author_id, author_name)
			VALUES (?, ?, ?, ?, ?)`
	err := tx.Raw(sqlStrCreateComment, comment.CommentID, comment.Content, comment.PostID, comment.AuthorID, comment.AuthorName).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	sqlStrCreateRelation := `
		INSERT INTO comment_relation (post_id, comment_id, parent_id, reply_id) 
			VALUES (?, ?, ?, ?)`
	err = tx.Raw(sqlStrCreateRelation, comment.PostID, comment.CommentID, parentID, replyID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// UpdateComment 更新评论
func (cd *CommentDao) UpdateComment(ctx context.Context, commentID int64, content string) error {
	sqlStr := `
		UPDATE comment
		SET content = ?
		WHERE comment_id = ?`
	return cd.WithContext(ctx).Exec(sqlStr, content, commentID).Error
}

// DeleteComment 删除评论
func (cd *CommentDao) DeleteComment(ctx context.Context, commentID int64) error {
	sqlStrDeleteComment := `
		UPDATE comment
		SET delete_time = ?
		WHERE comment_id = ?`
	sqlStrDeleteRelation := `
		UPDATE comment_relation
		SET delete_time = ?
		WHERE comment_id = ? OR parent_id = ? OR reply_id = ?`
	tx := cd.Begin().WithContext(ctx)
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
func (cd *CommentDao) GetCommentCount(ctx context.Context, postID int64) (int64, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*) FROM comment
		WHERE post_id = ? AND status = 1 AND delete_time = 0`
	err := cd.WithContext(ctx).Raw(sqlStr, postID).Scan(&count).Error
	return count, err
}

// GetTopCommentCount 获取顶级评论数量
func (cd *CommentDao) GetTopCommentCount(ctx context.Context, postID int64) (int64, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*) FROM comment
		INNER JOIN comment_relation ON comment.comment_id = comment_relation.comment_id
		WHERE comment.post_id = ? AND status = 1 AND comment.delete_time = 0 AND parent_id = 0 AND comment_relation.delete_time = 0`
	err := cd.WithContext(ctx).Raw(sqlStr, postID).Scan(&count).Error
	return count, err
}

func (cd *CommentDao) GetSubCommentCount(ctx context.Context, parentID int64) (int64, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*) FROM comment
		INNER JOIN comment_relation ON comment.comment_id = comment_relation.parent_id
		WHERE parent_id = ? AND status = 1 AND comment.delete_time = 0 AND comment_relation.delete_time = 0`
	err := cd.WithContext(ctx).Raw(sqlStr, parentID).Scan(&count).Error
	return count, err
}

// GetCommentCountByUserID 获取用户评论数量
func (cd *CommentDao) GetCommentCountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*) FROM comment
		WHERE author_id = ? AND status = 1 AND delete_time = 0`
	err := cd.WithContext(ctx).Raw(sqlStr, userID).Scan(&count).Error
	return count, err
}
