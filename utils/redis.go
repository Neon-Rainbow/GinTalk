package utils

import "fmt"

const (
	UserTokenKeyTemplate = "user:token:%v"
)

func GenerateRedisKey(template string, param ...interface{}) string {
	return fmt.Sprintf(template, param...)
}
