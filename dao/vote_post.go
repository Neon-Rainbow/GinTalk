package dao

import (
	"GinTalk/DTO"
	"GinTalk/model"
	"context"
	"gorm.io/gorm"
	"time"
)

var _ PostVoteDaoInterface = (*PostVoteDao)(nil)

type PostVoteDaoInterface interface {
	AddPostVote(ctx context.Context, postID int64, userID int64) error
	RevokePostVote(ctx context.Context, postID int64, userID int64) error

	// AddContentVoteUp 帖子赞成票数加一
	AddContentVoteUp(ctx context.Context, postID int64) error
	// SubContentVoteUp 帖子赞成票数减一
	SubContentVoteUp(ctx context.Context, postID int64) error

	// GetUserVoteList 获取用户投票过的帖子
	GetUserVoteList(ctx context.Context, userID int64, pageNum int, pageSize int) ([]int64, error)

	// GetPostVoteCount 获取帖子的赞成票数
	GetPostVoteCount(ctx context.Context, postID int64) (*DTO.PostVoteCounts, error)

	// GetBatchPostVoteCount 批量获取帖子的赞成票数
	GetBatchPostVoteCount(ctx context.Context, postID []int64) ([]DTO.PostVoteCounts, error)

	// CheckUserVoted 检查用户是否投过票
	CheckUserVoted(ctx context.Context, postID []int64, userID int64) ([]DTO.UserVotePostRelationsDTO, error)

	CheckUserVotedPost(ctx context.Context, postID int64, userID int64) (bool, error)

	// GetPostVoteDetail 获取帖子的投票详情
	GetPostVoteDetail(ctx context.Context, postID int64, pageNum int, pageSize int) ([]DTO.UserVotePostDetailDTO, error)

	// GetPostCreateTime 获取帖子的创建时间
	GetPostCreateTime(ctx context.Context, postID int64) (time.Time, error)
}

type PostVoteDao struct {
	*gorm.DB
}

func (p *PostVoteDao) AddPostVote(ctx context.Context, postID int64, userID int64) error {
	vote := model.VotePost{
		PostID: postID,
		UserID: userID,
		Vote:   1,
	}
	sqlStr := `
		INSERT INTO vote_post (post_id, user_id, vote)
		VALUES (?, ?, ?)	
`
	return p.DB.WithContext(ctx).Exec(sqlStr, vote.PostID, vote.UserID, vote.Vote).Error
}

func (p *PostVoteDao) RevokePostVote(ctx context.Context, postID int64, userID int64) error {
	sqlStr := `
		DELETE FROM vote_post
		WHERE post_id = ? AND user_id = ?`
	return p.DB.WithContext(ctx).Exec(sqlStr, postID, userID).Error
}

func (p *PostVoteDao) AddContentVoteUp(ctx context.Context, postID int64) error {
	sqlStr := `
		UPDATE content_votes
		SET vote = vote + 1
		WHERE post_id = ? AND delete_time = 0`
	return p.DB.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (p *PostVoteDao) SubContentVoteUp(ctx context.Context, postID int64) error {
	sqlStr := `
		UPDATE content_votes
		SET vote = vote - 1
		WHERE post_id = ? AND delete_time = 0`
	return p.DB.WithContext(ctx).Exec(sqlStr, postID).Error
}

func (p *PostVoteDao) GetUserVoteList(ctx context.Context, userID int64, pageNum int, pageSize int) ([]int64, error) {
	var voteRecord []int64
	sqlStr := `
		SELECT post_id
		FROM vote_post
		WHERE user_id = ?
		LIMIT ? OFFSET ?`
	err := p.DB.WithContext(ctx).Raw(sqlStr, userID, pageSize, (pageNum-1)*pageSize).Scan(&voteRecord).Error
	return voteRecord, err
}

func (p *PostVoteDao) GetPostVoteCount(ctx context.Context, postID int64) (*DTO.PostVoteCounts, error) {
	var voteCount DTO.PostVoteCounts
	sqlStr := `
		SELECT post_id, vote
		FROM content_votes
		WHERE post_id = ? AND delete_time = 0`
	err := p.DB.WithContext(ctx).Raw(sqlStr, postID).Scan(&voteCount).Error
	return &voteCount, err
}

func (p *PostVoteDao) GetBatchPostVoteCount(ctx context.Context, postIDs []int64) ([]DTO.PostVoteCounts, error) {
	var voteCount []DTO.PostVoteCounts
	sqlStr := `
		SELECT post_id, vote
		FROM content_votes
		WHERE post_id IN (?) AND delete_time = 0`
	err := p.DB.WithContext(ctx).Raw(sqlStr, postIDs).Scan(&voteCount).Error
	return voteCount, err
}

func (p *PostVoteDao) CheckUserVoted(ctx context.Context, postIDs []int64, userID int64) ([]DTO.UserVotePostRelationsDTO, error) {
	var votes []DTO.UserVotePostRelationsDTO
	sqlStr := `
		SELECT post_id, vote
		FROM vote_post
		WHERE post_id IN (?) AND user_id = ? AND delete_time = 0`
	err := p.DB.WithContext(ctx).Raw(sqlStr, postIDs, userID).Scan(&votes).Error
	return votes, err
}

func (p *PostVoteDao) GetPostVoteDetail(ctx context.Context, postID int64, pageNum int, pageSize int) ([]DTO.UserVotePostDetailDTO, error) {
	var votes []DTO.UserVotePostDetailDTO
	sqlStr := `
		SELECT user.user_id, post_id, vote, username
		FROM vote_post
		INNER JOIN user ON vote_post.user_id = user.user_id
		WHERE post_id = ? AND vote_post.delete_time = 0 AND user.delete_time = 0
		LIMIT ? OFFSET ?`
	err := p.DB.WithContext(ctx).Raw(sqlStr, postID, pageSize, (pageNum-1)*pageSize).Scan(&votes).Error
	return votes, err
}

func (p *PostVoteDao) GetPostCreateTime(ctx context.Context, postID int64) (time.Time, error) {
	var createTime time.Time
	sqlStr := `
		SELECT create_time
		FROM post
		WHERE post_id = ?`
	err := p.DB.WithContext(ctx).Raw(sqlStr, postID).Scan(&createTime).Error
	return createTime, err
}

func (p *PostVoteDao) CheckUserVotedPost(ctx context.Context, postID int64, userID int64) (bool, error) {
	var count int64
	sqlStr := `
		SELECT COUNT(*)
		FROM vote_post
		WHERE post_id = ? AND user_id = ? AND delete_time = 0`
	err := p.DB.WithContext(ctx).Raw(sqlStr, postID, userID).Scan(&count).Error
	return count > 0, err
}

func NewPostVoteDao(db *gorm.DB) *PostVoteDao {
	return &PostVoteDao{DB: db}
}
