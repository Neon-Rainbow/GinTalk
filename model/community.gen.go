// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameCommunity = "community"

// Community 社区表：存储社区信息
type Community struct {
	ID            int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CommunityID   int32     `gorm:"column:community_id;not null" json:"community_id"`
	CommunityName string    `gorm:"column:community_name;not null" json:"community_name"`
	Introduction  string    `gorm:"column:introduction;not null" json:"introduction"`
	CreateTime    time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
	DeleteTime    int       `gorm:"column:delete_time" json:"delete_time"`
}

// TableName Community's table name
func (*Community) TableName() string {
	return TableNameCommunity
}
