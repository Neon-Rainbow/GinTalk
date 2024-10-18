package DTO

type PostDetail struct {
	PostID        uint   `json:"post_id" db:"post_id"`
	Title         string `json:"title" db:"title"`
	Content       string `json:"content" db:"content"`
	AuthorId      uint   `json:"author_id" db:"author_id"`
	Username      string `json:"author_name" db:"username"`
	CommunityID   int64  `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
	Status        int32  `json:"status" db:"status"`
}
