package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/user/request"
	"takeout/internal/api/user/response"
	"takeout/internal/service"
)

type WxUserController struct {
	service service.IWxUserService
}

func NewWxUserController(service service.IWxUserService) *WxUserController {
	return &WxUserController{service: service}
}

// Login @Login 微信小程序登录
// @Tags WxUser
// @Security JWTAuth
// @Produce json
// @Param data body request.WxUserLoginDTO true "微信授权码"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/user/login [post]
func (c *WxUserController) Login(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		err  error
		dto  request.WxUserLoginDTO
		vo   response.WxUserVO
	)
	if err = ctx.ShouldBind(&dto); err != nil {
		global.Log.Debug("bind param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}
	if vo, err = c.service.Login(ctx, dto); err != nil {
		code = e.ERROR
		global.Log.Debug("login error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: vo,
		Msg:  e.GetMsg(code),
	})

}

// Logout @Logout 微信小程序退出登录
// @Tags WxUser
// @Security JWTAuth
// @Produce json
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/user/logout [post]
func (c *WxUserController) Logout(ctx *gin.Context) {
	var code = e.SUCCESS

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}
