package DTO

type VoteDTO struct {
	PostID int64 `json:"post_id" binding:"required"`
	UserID int64 `json:"user_id" binding:"required"`
	Vote   int   `json:"vote"`
}

type MyVoteListDTO struct {
	UserID   int64 `json:"user_id" form:"user_id"`
	PageNum  int   `json:"page_num" form:"page_num"`
	PageSize int   `json:"page_size" form:"page_size"`
}

type CreateVoteDTO struct {
	UserID int64   `json:"user_id" binding:"required"`
	PostID []int64 `json:"post_id" binding:"required"`
}

type VoteCountDTO struct {
	PostID int64 `json:"post_id" binding:"required"`
}
