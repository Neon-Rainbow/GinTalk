// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameComment = "comment"

// Comment mapped from table <comment>
type Comment struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增主键，唯一标识每条评论记录" json:"id"`                // 自增主键，唯一标识每条评论记录
	CommentID  int64          `gorm:"column:comment_id;not null;comment:评论ID，用于业务中的评论唯一标识" json:"comment_id"`                   // 评论ID，用于业务中的评论唯一标识
	Content    string         `gorm:"column:content;not null;comment:评论内容" json:"content"`                                      // 评论内容
	PostID     int64          `gorm:"column:post_id;not null;comment:评论所属的帖子ID" json:"post_id"`                                 // 评论所属的帖子ID
	AuthorID   int64          `gorm:"column:author_id;not null;comment:评论作者的用户ID" json:"author_id"`                             // 评论作者的用户ID
	AuthorName string         `gorm:"column:author_name;not null;comment:评论时的用户的名字" json:"author_name"`                         // 评论时的用户的名字
	ParentID   int64          `gorm:"column:parent_id;not null;comment:该评论回复的评论ID，为0表示原生评论,即第一层的评论，不为0表示回复评论" json:"parent_id"` // 该评论回复的评论ID，为0表示原生评论,即第一层的评论，不为0表示回复评论
	ReplyID    int64          `gorm:"column:reply_id;not null;comment:父评论ID, 为0表示原生评论，不为0表示回复评论" json:"reply_id"`               // 父评论ID, 为0表示原生评论，不为0表示回复评论
	Status     int32          `gorm:"column:status;not null;default:1;comment:评论状态：1-正常，0-删除" json:"status"`                    // 评论状态：1-正常，0-删除
	CreateTime time.Time      `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:评论创建时间，默认当前时间" json:"create_time"`    // 评论创建时间，默认当前时间
	UpdateTime time.Time      `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:评论更新时间，每次更新时自动修改" json:"update_time"` // 评论更新时间，每次更新时自动修改
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;comment:逻辑删除时间，NULL表示未删除" json:"delete_time"`                           // 逻辑删除时间，NULL表示未删除
}

// TableName Comment's table name
func (*Comment) TableName() string {
	return TableNameComment
}
