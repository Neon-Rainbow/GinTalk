package DTO

type Comment struct {
	CommentID  int64  `json:"comment_id" db:"comment_id"`
	PostID     int64  `json:"post_id" db:"post_id"`
	AuthorID   int64  `json:"author_id" db:"author_id"`
	AuthorName string `json:"author_name" db:"author_name"`
	Content    string `json:"content" db:"content"`
}

type CommentRequest struct {
	*Comment
	ReplyID  int64 `json:"reply_id" db:"reply_id"`
	ParentID int64 `json:"parent_id" db:"parent_id"`
}