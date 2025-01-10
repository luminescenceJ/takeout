package dao

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"takeout/internal/model"
)

type WxUserDao struct {
	db *gorm.DB
}

func (w WxUserDao) RegisterUser(ctx *gin.Context, user model.User) error {
	return w.db.WithContext(ctx).Create(&user).Error
}

func (w WxUserDao) GetUserByOpenID(ctx *gin.Context, openID string) (model.User, bool, error) {
	var (
		user model.User
		err  error
	)

	if err = w.db.WithContext(ctx).Where("openid = ?", openID).First(&user).Error; err != nil {
		// 如果是记录未找到的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, false, nil
		}

		// 如果是其他数据库查询错误
		if err != nil {
			return model.User{}, false, err
		}
	}
	return user, true, nil
}

func NewWxUserDao(db *gorm.DB) *WxUserDao {
	return &WxUserDao{db: db}
}
