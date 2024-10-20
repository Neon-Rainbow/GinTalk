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
