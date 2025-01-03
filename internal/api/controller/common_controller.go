package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"takeout/common"
	"takeout/common/e"
	"takeout/common/utils"
	"takeout/global"
)

type CommonController struct {
}

func (c *CommonController) Upload(ctx *gin.Context) {
	code := e.SUCCESS
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}
	uuid := uuid.New()
	imageName := uuid.String() + file.Filename
	imagePath, err := utils.LocalOss(imageName, file)
	//imagePath, err := utils.AliyunOss(imageName, file)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("AliyunOss upload failed", "err", err.Error())
	}

	ctx.JSON(http.StatusOK, common.Result{Code: code, Data: imagePath, Msg: e.GetMsg(code)})

}
