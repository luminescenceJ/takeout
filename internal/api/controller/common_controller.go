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

// Upload 上传文件
// @Summary 上传文件到服务器
// @Description 接收一个文件，并将其上传到本地存储或云存储（如阿里云OSS）
// @Tags FileUpload
// @Accept multipart/form-data
// @Security JWTAuth
// @Produce json
// @Param file formData file true "文件"  // file 参数表示文件
// @Success 200 {object} common.Result{code=int, data=string, msg=string} "上传成功"
// @Failure 400 {object} common.Result{code=int, data=string, msg=string} "上传失败"
// @Router /admin/common/upload [post]
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
