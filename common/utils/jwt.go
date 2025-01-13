package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// CustomPayload 自定义载荷继承原有接口并附带自己的字段
type CustomPayload struct {
	UserId     uint64
	GrantScope string
	jwt.RegisteredClaims
}

// GenerateToken 生成Token uid 用户id grantScope 签发对象  secret 加盐
func GenerateToken(userId uint64, grantScope string, secret string) (string, error) {
	claims := CustomPayload{
		UserId:     userId,
		GrantScope: grantScope,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                      //签发者
			Subject:   grantScope,                                         //签发对象
			Audience:  jwt.ClaimStrings{"PC", "Wechat_Program"},           //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), //过期时间
			NotBefore: jwt.NewNumericDate(time.Now()),                     //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     //签发时间
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	return token, err
}

func ParseToken(token string, secret string) (*CustomPayload, error) {
	parseToken, err := jwt.ParseWithClaims(token, &CustomPayload{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parseToken.Claims.(*CustomPayload); ok && parseToken.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
