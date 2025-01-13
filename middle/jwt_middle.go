package middle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/common/enum"
	"takeout/common/utils"
	"takeout/global"
)

func VerifiyJWTAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(global.Config.Jwt.Admin.Name)
		// 解析获取用户载荷信息
		payLoad, err := utils.ParseToken(token, global.Config.Jwt.Admin.Secret)
		if err != nil {
			code := e.UNKNOW_IDENTITY
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 在上下文设置载荷信息
		c.Set(enum.CurrentId, payLoad.UserId)
		c.Set(enum.CurrentName, payLoad.GrantScope)
		c.Next()
	}
}

func VerifiyJWTUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(global.Config.Jwt.User.Name)
		// 解析获取用户载荷信息
		payLoad, err := utils.ParseToken(token, global.Config.Jwt.User.Secret)
		if err != nil {
			code := e.UNKNOW_IDENTITY
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 在上下文设置载荷信息
		c.Set(enum.CurrentId, payLoad.UserId)
		c.Set(enum.CurrentName, payLoad.GrantScope)
		c.Next()
	}
}
