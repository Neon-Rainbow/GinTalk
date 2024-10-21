package dao

import (
	"GinTalk/DTO"
	"GinTalk/model"
	"context"
	"fmt"
	"gorm.io/gorm"
)

var _ PostDaoInterface = (*PostDao)(nil)

type PostDaoInterface interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPostList(ctx context.Context, pageNum int, pageSize int) ([]*DTO.PostSummary, error)
	GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDetail, error)

	// UpdatePost 更新帖子
	// 参数post中只需要有PostID, Title, Summary, Content四个字段
	UpdatePost(ctx context.Context, post *model.Post) error
}

type PostDao struct {
	*gorm.DB
}

func NewPostDao(db *gorm.DB) PostDaoInterface {
	return &PostDao{DB: db}
}

func (pd *PostDao) CreatePost(ctx context.Context, post *model.Post) error {
	if post.PostID == 0 {
		return fmt.Errorf("postID 不能为空")
	}
	if post.Title == "" {
		return fmt.Errorf("标题不能为空")
	}
	if post.Content == "" {
		return fmt.Errorf("内容不能为空")
	}
	if post.Summary == "" {
		return fmt.Errorf("摘要不能为空")
	}
	if post.AuthorID == 0 {
		return fmt.Errorf("作者ID不能为空")
	}
	if post.CommunityID == 0 {
		return fmt.Errorf("社区ID不能为空")
	}

	sqlStr1 := `INSERT INTO post (post_id, title,summary, content, author_id, community_id) VALUES (?, ?, ?,?, ?, ?)`
	sqlStr2 := `INSERT INTO content_votes (post_id) VALUES (?)`

	tx := pd.WithContext(ctx).Begin()
	err := tx.WithContext(ctx).Exec(sqlStr1, post.PostID, post.Title, post.Summary, post.Content, post.AuthorID, post.CommunityID).Error
	if err != nil {
		tx.Rollback()
	}
	err = tx.WithContext(ctx).Exec(sqlStr2, post.PostID).Error
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (pd *PostDao) GetPostList(ctx context.Context, pageNum int, pageSize int) ([]*DTO.PostSummary, error) {
	sqlStr := `SELECT 
                    post.post_id,
                    post.title,
                    post.summary,
                    post.author_id,
                    user.username,
                    post.community_id,
                    community.community_name,
                    post.status 
                FROM 
                    post
                INNER JOIN 
                    community ON community.community_id = post.community_id
                INNER JOIN 
                    user ON user.user_id = post.author_id
                WHERE 
                    post.delete_time = 0 
                LIMIT ? OFFSET ?`

	var posts []*DTO.PostSummary
	err := pd.WithContext(ctx).Raw(sqlStr, pageSize, (pageNum-1)*pageSize).Scan(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (pd *PostDao) GetPostDetail(ctx context.Context, postID int64) (*DTO.PostDetail, error) {
	sqlStr := `SELECT 
					post.post_id,
					post.title,
					post.content,
					post.author_id,
					user.username,
					post.community_id,
					community.community_name,
					post.status 
				FROM 
					post
				INNER JOIN 
					community ON community.community_id = post.community_id
				INNER JOIN 
					user ON user.user_id = post.author_id
				WHERE 
					post.post_id = ? 
					AND post.delete_time = 0`

	var post DTO.PostDetail
	err := pd.WithContext(ctx).Raw(sqlStr, postID).Scan(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (pd *PostDao) UpdatePost(ctx context.Context, post *model.Post) error {
	if post.PostID == 0 {
		return fmt.Errorf("postID 不能为空")
	}
	if post.Title == "" {
		return fmt.Errorf("标题不能为空")
	}
	if post.Content == "" {
		return fmt.Errorf("内容不能为空")
	}
	if post.Summary == "" {
		return fmt.Errorf("摘要不能为空")
	}
	sqlStr := `UPDATE post SET title = ?, summary = ?,content = ? WHERE post_id = ?`
	err := pd.WithContext(ctx).Exec(sqlStr, post.Title, post.Summary, post.Content, post.PostID).Error
	if err != nil {
		return err
	}
	return nil
}
