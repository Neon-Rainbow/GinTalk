package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CommentID uint64 `gorm:"uniqueIndex;not null"`
	Content   string `gorm:"type:text;not null"`
	PostID    uint64 `gorm:"index;not null"`
	AuthorID  uint64 `gorm:"index;not null"`
	ParentID  uint64 `gorm:"default:0;not null"`
	Status    uint8  `gorm:"default:1;not null"`
}
