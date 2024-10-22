package DTO

type PostDetail struct {
	PostID        int64  `json:"post_id,omitempty" db:"post_id"`
	Title         string `json:"title,omitempty" db:"title"`
	Content       string `json:"content,omitempty" db:"content"`
	AuthorId      int64  `json:"author_id,omitempty" db:"author_id"`
	Username      string `json:"author_name,omitempty" db:"username"`
	CommunityID   int64  `json:"community_id,omitempty" db:"community_id"`
	CommunityName string `json:"community_name,omitempty" db:"community_name"`
	Status        int32  `json:"status,omitempty" db:"status"`
}

// PostSummary 帖子摘要
// 用于帖子列表展示
type PostSummary struct {
	PostID        int64  `json:"post_id,omitempty" db:"post_id"`
	Title         string `json:"title,omitempty" db:"title"`
	Summary       string `json:"summary,omitempty" db:"summary"`
	AuthorId      int64  `json:"author_id,omitempty" db:"author_id"`
	Username      string `json:"author_name,omitempty" db:"username"`
	CommunityID   int64  `json:"community_id,omitempty" db:"community_id"`
	CommunityName string `json:"community_name,omitempty" db:"community_name"`
}

// PostVoteCounts 帖子投票内容
// 用于获取帖子的投票内容
type PostVoteCounts struct {
	PostID int64 `json:"post_id,omitempty" db:"post_id"`
	Vote   int64 `json:"vote" db:"vote"`
}

type UserVotePostRelationsDTO struct {
	UserID int64 `json:"user_id"`
	PostID int64 `json:"post_id,omitempty"`
	Vote   int   `json:"vote"`
}
