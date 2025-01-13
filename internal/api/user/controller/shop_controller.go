package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/service"
)

type ShopController struct {
	shopService *service.ShopService
}

// GetShopStatus @GetShopStatus 获取营业状态
// @Tags UserShop
// @Security JWTAuth
// @Produce json
// @Success 200 {object} common.Result "success"
// @Router /user/shop/status [get]
func (s *ShopController) GetShopStatus(ctx *gin.Context) {
	status, err := s.shopService.GetShopStatus()
	if err != nil {
		global.Log.Debug("GetShopStatus error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{})
		return
	}
	statusInt, _ := strconv.Atoi(status)
	ctx.JSON(http.StatusOK, common.Result{
		Code: e.SUCCESS,
		Data: statusInt,
		Msg:  e.GetMsg(e.SUCCESS),
	})
}
