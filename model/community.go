package model

import "gorm.io/gorm"

type Community struct {
	gorm.Model
	CommunityID   uint   `gorm:"uniqueIndex;not null"`
	CommunityName string `gorm:"type:varchar(128);uniqueIndex;not null"`
	Introduction  string `gorm:"type:varchar(256);not null"`
}
