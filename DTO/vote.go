package DTO

type VoteDTO struct {
	ID      int64 `json:"id" binding:"required"`       // 帖子ID或评论ID
	VoteFor int   `json:"vote_for" binding:"required"` // 1: 帖子 2: 评论
	UserID  int64 `json:"user_id" binding:"required"`  // 用户ID
	Vote    int   `json:"vote"`                        // -1: 踩 0: 取消 1: 赞
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
