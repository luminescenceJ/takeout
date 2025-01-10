package repository

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/model"
)

type WxUserRepo interface {
	GetUserByOpenID(ctx *gin.Context, openID string) (model.User, bool, error)
	RegisterUser(ctx *gin.Context, user model.User) error
}
