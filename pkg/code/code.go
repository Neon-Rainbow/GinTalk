package code

type RespCode uint

const (
	Success RespCode = 1000 + iota
	InvalidParam
	InvalidPassword
	InvalidToken
	InvalidAuth
	ServerError
	UserNotExist
	PasswordError
	UserRefreshTokenError
	TimeOut
)

var codeMsg = map[RespCode]string{
	Success:               "success",
	InvalidParam:          "请求参数错误",
	InvalidPassword:       "密码错误",
	InvalidToken:          "无效的token",
	InvalidAuth:           "无效的授权",
	ServerError:           "服务器错误",
	UserNotExist:          "用户不存在",
	PasswordError:         "密码错误",
	UserRefreshTokenError: "刷新token错误",
	TimeOut:               "超时",
}

func (c RespCode) GetMsg() string {
	msg, ok := codeMsg[c]
	if !ok {
		return codeMsg[ServerError]
	}
	return msg
}
