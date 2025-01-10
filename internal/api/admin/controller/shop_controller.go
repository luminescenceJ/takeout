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

// SetShopStatus @SetShopStatus 设置营业状态
// @Tags shop
// @Security JWTAuth
// @Produce json
// @Param status path string true "店铺营业状态：1为营业，0为打烊"
// @Success 200 {object} common.Result "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /admin/shop/{status} [put]
func (s *ShopController) SetShopStatus(ctx *gin.Context) {
	var statusMap = map[string]string{"0": "休息", "1": "营业"}
	status := ctx.Param("status")
	if status == "" {
		global.Log.Debug("status parameter is empty")
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}
	if err := s.shopService.SetShopStatus(status); err != nil {
		global.Log.Debug("SetShopStatus error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{})
		return
	}
	global.Log.Info("设置店铺营业状态: ", statusMap[status])
	ctx.JSON(http.StatusOK, common.Result{
		Code: e.SUCCESS,
		Msg:  e.GetMsg(e.SUCCESS),
	})
}

// GetShopStatus @GetShopStatus 获取营业状态
// @Tags shop
// @Security JWTAuth
// @Produce json
// @Success 200 {object} common.Result "success"
// @Router /admin/shop/status [get]
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
