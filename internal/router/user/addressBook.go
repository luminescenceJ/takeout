package user

import (
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/internal/api/user/controller"
	"takeout/internal/service"
	"takeout/middle"
	"takeout/repository/dao"
)

type AddressBookRouter struct{}

func (dr *AddressBookRouter) InitApiRouter(parent *gin.RouterGroup) {
	privateRouter := parent.Group("addressBook") // 私有路由使用jwt验证
	privateRouter.Use(middle.VerifiyJWTUser())
	// 依赖注入
	addressCtrl := controller.NewAddressBookController(
		service.NewAddressBookService(dao.NewAddressBookDao(global.DB)),
	)
	{
		privateRouter.POST("/", addressCtrl.Create)
		privateRouter.PUT("/", addressCtrl.Update)
		privateRouter.DELETE("/", addressCtrl.DeleteById)
		privateRouter.GET("/:id", addressCtrl.GetById)

		privateRouter.GET("/list", addressCtrl.GetCurAddress)
		privateRouter.GET("/default", addressCtrl.GetDefaultAddress)
		privateRouter.PUT("/default", addressCtrl.SetDefaultAddress)

	}
}
