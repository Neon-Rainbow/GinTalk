package pkg

import (
	"GinTalk/settings"
	"crypto/md5"
	"encoding/hex"
)

// EncryptPassword 用于加密密码
func EncryptPassword(password string) string {
	var secret = settings.GetConfig().PasswordSecret
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
