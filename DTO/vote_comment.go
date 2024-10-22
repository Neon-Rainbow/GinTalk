package DTO

type VoteComment struct {
	UserID    int64 `json:"user_id" form:"user_id"`
	CommentID int64 `json:"comment_id" form:"comment_id"`
}
