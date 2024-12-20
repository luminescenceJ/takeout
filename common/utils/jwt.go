package utils

import "github.com/golang-jwt/jwt/v5"

// CustomPayload 自定义载荷继承原有接口并附带自己的字段
type CustomPayload struct {
	UserId     uint64
	GrantScope string
	jwt.RegisteredClaims
}
