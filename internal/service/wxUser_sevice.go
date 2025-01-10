package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"takeout/common/utils"
	"takeout/global"
	"takeout/internal/api/user/request"
	"takeout/internal/api/user/response"
	"takeout/internal/model"
	"takeout/repository"
)

type IWxUserService interface {
	Login(ctx *gin.Context, request request.WxUserLoginDTO) (response.WxUserVO, error)
}

type WxUserService struct {
	repo repository.WxUserRepo
}

func (ws WxUserService) Login(ctx *gin.Context, request request.WxUserLoginDTO) (response.WxUserVO, error) {
	var (
		openID    = utils.GetOpenID(request.Code)
		user      model.User
		err       error
		existUser bool
		jwtToken  string
	)
	// 获取微信openid
	if openID == "" {
		return response.WxUserVO{}, errors.New("openId is empty")
	}
	// 判断用户是否存在于数据库 存在就返回token 不存在就自动注册
	if user, existUser, err = ws.repo.GetUserByOpenID(ctx, openID); err != nil {
		return response.WxUserVO{}, err
	}
	//用户未注册
	if !existUser {
		user.OpenId = openID
		if err = ws.repo.RegisterUser(ctx, user); err != nil {
			//注册失败
			return response.WxUserVO{}, err
		}
	}
	// 用户存在，分发jwt令牌
	jwtConfig := global.Config.Jwt.Admin
	if jwtToken, err = utils.GenerateToken(uint64(user.ID), jwtConfig.Name, jwtConfig.Secret); err != nil {
		return response.WxUserVO{}, err
	}
	// 返回登录结果
	res := response.WxUserVO{
		Id:     int64(user.ID),
		Token:  jwtToken,
		Openid: openID,
	}
	return res, nil
}

func NewWxUserService(repo repository.WxUserRepo) IWxUserService {
	return &WxUserService{repo: repo}
}
