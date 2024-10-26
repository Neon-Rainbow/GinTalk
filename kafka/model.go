package kafka

type Vote struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
	Vote   int    `json:"vote"`
}
