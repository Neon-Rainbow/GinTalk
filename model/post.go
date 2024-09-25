package model

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	PostID      uint64 `gorm:"uniqueIndex;not null"`
	Title       string `gorm:"type:varchar(128);not null"`
	Content     string `gorm:"type:varchar(8192);not null"`
	AuthorID    uint64 `gorm:"index;not null"`
	CommunityID uint64 `gorm:"index;not null"`
	Status      int8   `gorm:"default:1;not null"`
}
