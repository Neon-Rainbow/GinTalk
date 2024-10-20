package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"sync"
	"time"
)

var mySecret = "111"

const (
	// AccessTokenName 是访问令牌的key
	AccessTokenName = "access"
	// RefreshTokenName 是刷新令牌的key
	RefreshTokenName = "refresh"
)

type MyClaims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateToken 生成token
func GenerateToken[T int64 | string | uint](userID T, username string) (accessToken string, refreshToken string, err error) {
	var int64UserID int64
	switch v := any(userID).(type) {
	case uint:
		int64UserID = int64(v)
	case int64:
		int64UserID = v
	case string:
		// 尝试将 string 转为 uint
		parsedID, err := strconv.ParseUint(v, 10, 32) // 假设 uint 是 32 位
		if err != nil {
			return "", "", fmt.Errorf("invalid userID format, could not convert to uint: %v", err)
		}
		int64UserID = int64(parsedID)
	default:
		return "", "", fmt.Errorf("unsupported userID type")
	}

	f := func(userID int64, username string, tokenType string, validTime time.Duration) (string, error) {
		c := MyClaims{
			UserID:    userID,
			Username:  username,
			TokenType: tokenType,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(validTime)),
				Issuer:    "水告木南",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
		return token.SignedString([]byte(mySecret))
	}

	var wg sync.WaitGroup
	errorChannel := make(chan error)
	wg.Add(2)
	go func() {
		defer wg.Done()
		accessToken, err = f(int64UserID, username, AccessTokenName, time.Hour*24*7)
		if err != nil {
			errorChannel <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		refreshToken, err = f(int64UserID, username, RefreshTokenName, time.Hour*24*7)
		if err != nil {
			errorChannel <- err
			return
		}
	}()

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	for err = range errorChannel {
		if err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*MyClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}
