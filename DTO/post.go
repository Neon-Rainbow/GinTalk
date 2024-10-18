package DTO

type PostDTO struct {
	PostID      uint   `json:"post_id" db:"post_id"`
	Title       string `json:"title" db:"title"`
	Content     string `json:"content" db:"content"`
	AuthorId    uint   `json:"author_id" db:"author_id"`
	CommunityID int64  `json:"community_id" db:"community_id"`
	Status      int32  `json:"status" db:"status"`
}
