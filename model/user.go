package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(64);not null"`
	Email    string `gorm:"type:varchar(64)"`
	Gender   int8   `gorm:"default:0;not null"`
}
