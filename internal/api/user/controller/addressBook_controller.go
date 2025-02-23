package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"takeout/common"
	"takeout/common/e"
	"takeout/global"
	"takeout/internal/api/user/request"
	"takeout/internal/model"
	"takeout/internal/service"
)

type AddressBookController struct {
	service service.IAddressBookService
}

func NewAddressBookController(service service.IAddressBookService) *AddressBookController {
	return &AddressBookController{service: service}
}

// Create @Create 新增地址
// @Tags C端-地址簿接口
// @Security JWTAuth
// @Produce json
// @Param address body request.AddressBookDTO true "查询信息"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/addressBook [post]
func (ac *AddressBookController) Create(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		address request.AddressBookDTO
		err     error
	)

	if err = ctx.ShouldBindJSON(&address); err != nil {
		global.Log.Debug("C端-地址簿接口 Create param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = ac.service.CreateAddressBook(ctx, address); err != nil {
		code = e.ERROR
		global.Log.Debug("CreateAddressBook error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})

}

// Update @Update 根据id修改地址
// @Tags C端-地址簿接口
// @Security JWTAuth
// @Produce json
// @Param address body model.AddressBook true "修改信息"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/addressBook [put]
func (ac *AddressBookController) Update(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		address request.AddressBookDTO
		err     error
	)

	if err = ctx.ShouldBind(&address); err != nil {
		global.Log.Debug("C端-地址簿接口 Update param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = ac.service.UpdateAddressBook(ctx, address); err != nil {
		code = e.ERROR
		global.Log.Debug("UpdateAddressBook error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// DeleteById @DeleteById 根据id删除地址
// @Tags C端-地址簿接口
// @Security JWTAuth
// @Produce json
// @Param id query string true "地址id"
// @Success 200 {object} common.Result "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/addressBook [delete]
func (ac *AddressBookController) DeleteById(ctx *gin.Context) {
	var (
		code = e.SUCCESS
		id   uint64
		err  error
	)
	if id, err = strconv.ParseUint(ctx.Param("id"), 10, 64); err != nil {
		global.Log.Debug("C端-地址簿接口 DeleteById param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = ac.service.DeleteAddressBook(ctx, id); err != nil {
		code = e.ERROR
		global.Log.Debug("DeleteAddressBook error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

// GetById @GetById 根据id查询地址
// @Tags C端-地址簿接口
// @Security JWTAuth
// @Produce json
// @Param id query string true "地址id"
// @Success 200 {object} common.Result{Data=model.AddressBook} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/addressBook/{id} [get]
func (ac *AddressBookController) GetById(ctx *gin.Context) {
	var (
		code    = e.SUCCESS
		id      uint64
		address model.AddressBook
		err     error
	)
	if id, err = strconv.ParseUint(ctx.Param("id"), 10, 64); err != nil {
		global.Log.Debug("C端-地址簿接口 GetById param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if address, err = ac.service.GetAddressBook(ctx, id); err != nil {
		code = e.ERROR
		global.Log.Debug("GetAddressBook error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: address,
		Msg:  e.GetMsg(code),
	})
}

// GetCurAddress @GetCurAddress 查询当前登录用户的所有地址信息
// @Tags C端-地址簿接口
// @Security JWTAuth
// @Produce json
// @Success 200 {object} common.Result{Data=model.AddressBook} "success"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/addressBook/list [get]
func (ac *AddressBookController) GetCurAddress(ctx *gin.Context) {
	var (
		code        = e.SUCCESS
		addressList []model.AddressBook
		err         error
	)

	if addressList, err = ac.service.GetCurAddressBook(ctx); err != nil {
		code = e.ERROR
		global.Log.Debug("GetCurAddressBook error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: addressList,
		Msg:  e.GetMsg(code),
	})
}

// GetDefaultAddress @GetDefaultAddress 查询默认地址
// @Tags C端-地址簿接口
// @Security JWTAuth
// @Produce json
// @Success 200 {object} common.Result{Data=model.AddressBook} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/addressBook/default [get]
func (ac *AddressBookController) GetDefaultAddress(ctx *gin.Context) {
	var (
		code        = e.SUCCESS
		addressList model.AddressBook
		err         error
	)

	if addressList, err = ac.service.GetDefaultAddressBook(ctx); err != nil {
		code = e.ERROR
		global.Log.Debug("GetDefaultAddress error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Data: addressList,
		Msg:  e.GetMsg(code),
	})
}

// SetDefaultAddress @SetDefaultAddress 设置默认地址
// @Tags C端-地址簿接口
// @Security JWTAuth
// @Produce json
// @Param id body uint64 true "地址id"
// @Success 200 {object} common.Result{} "success"
// @Failure 400 {object} common.Result "Invalid request payload"
// @Failure 500 {object} common.Result "Internal Server Faliure"
// @Router /user/addressBook/default [put]
func (ac *AddressBookController) SetDefaultAddress(ctx *gin.Context) {
	type Request struct {
		ID uint64 `json:"id"`
	}

	var (
		code      = e.SUCCESS
		requestId Request
		err       error
	)

	if err = ctx.ShouldBindJSON(&requestId); err != nil {
		global.Log.Debug("C端-地址簿接口 SetDefaultAddress param error:", err.Error())
		ctx.JSON(http.StatusBadRequest, common.Result{})
		return
	}

	if err = ac.service.SetDefaultAddressBook(ctx, requestId.ID); err != nil {
		code = e.ERROR
		global.Log.Debug("SetDefaultAddressBook error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, common.Result{
			Code: code,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Result{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}
