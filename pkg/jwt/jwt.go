package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

var mySecret = "111"

type MyClaims struct {
	UserID    uint   `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateToken 生成token
func GenerateToken(userID uint) (accessToken string, refreshToken string, err error) {

	f := func(userID uint, tokenType string, validTime time.Duration) (string, error) {
		c := MyClaims{
			UserID:    userID,
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
		accessToken, err = f(userID, "access", time.Hour)
		if err != nil {
			errorChannel <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		refreshToken, err = f(userID, "refresh", time.Hour*24)
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
