package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"forum-gin/settings"
)

// EncryptPassword 用于加密密码
func EncryptPassword(password string) string {
	var secret = settings.Conf.PasswordSecret
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
