package DTO

type VoteDTO struct {
	ID      int64 `json:"id" binding:"required"`
	VoteFor int   `json:"vote_for" binding:"required"`
	UserID  int64 `json:"user_id" binding:"required"`
	Vote    int   `json:"vote"`
}

type MyVoteListDTO struct {
	UserID   int64 `json:"user_id" form:"user_id"`
	VoteFor  int   `json:"vote_for" form:"vote_for"`
	PageNum  int   `json:"page_num" form:"page_num"`
	PageSize int   `json:"page_size" form:"page_size"`
}

type CheckVoteListDTO struct {
	UserID  int64   `json:"user_id" form:"user_id" binding:"required"`
	IDs     []int64 `json:"id" form:"id" binding:"required"`
	VoteFor int     `json:"vote_for" form:"vote_for" binding:"required"`
}

type VoteCountDTO struct {
	ID      int64 `json:"id" form:"post_id" binding:"required"`
	VoteFor int   `json:"vote_for" form:"vote_for" binding:"required"`
}

type UserVoteDetailDTO struct {
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id,omitempty"`
	CommentID int64  `json:"comment_id,omitempty"`
	Username  string `json:"username"`
	Vote      int    `json:"vote"`
}
