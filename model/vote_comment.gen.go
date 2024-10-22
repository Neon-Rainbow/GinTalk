// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameVoteComment = "vote_comment"

// VoteComment 评论投票表：存储用户对评论的投票记录
type VoteComment struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增主键，唯一标识每条投票记录" json:"id"`                // 自增主键，唯一标识每条投票记录
	CommentID  int64     `gorm:"column:comment_id;not null;comment:投票所属的评论ID" json:"comment_id"`                           // 投票所属的评论ID
	UserID     int64     `gorm:"column:user_id;not null;comment:投票用户的用户ID" json:"user_id"`                                 // 投票用户的用户ID
	Vote       int32     `gorm:"column:vote;not null;comment:投票类型：1-赞" json:"vote"`                                        // 投票类型：1-赞
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:投票创建时间，默认当前时间" json:"create_time"`    // 投票创建时间，默认当前时间
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:投票更新时间，每次更新时自动修改" json:"update_time"` // 投票更新时间，每次更新时自动修改
	DeleteTime int       `gorm:"column:delete_time;comment:逻辑删除时间，NULL表示未删除" json:"delete_time"`                           // 逻辑删除时间，NULL表示未删除
}

// TableName VoteComment's table name
func (*VoteComment) TableName() string {
	return TableNameVoteComment
}
