// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameContentVote = "content_votes"

// ContentVote 帖子投票表：存储用户对帖子的投票记录
type ContentVote struct {
	PostID     int64     `gorm:"column:post_id;not null;comment:投票所属的帖子ID" json:"post_id"`                                 // 投票所属的帖子ID
	Count      int32     `gorm:"column:count;not null;comment:投票总数" json:"count"`                                          // 投票总数
	Vote       int32     `gorm:"column:vote;not null;comment:赞数" json:"vote"`                                              // 赞数
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:投票创建时间，默认当前时间" json:"create_time"`    // 投票创建时间，默认当前时间
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:投票更新时间，每次更新时自动修改" json:"update_time"` // 投票更新时间，每次更新时自动修改
	DeleteTime int       `gorm:"column:delete_time;comment:逻辑删除时间，NULL表示未删除" json:"delete_time"`                           // 逻辑删除时间，NULL表示未删除
}

// TableName ContentVote's table name
func (*ContentVote) TableName() string {
	return TableNameContentVote
}
