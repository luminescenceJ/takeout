package initialize

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	_ "takeout/docs" // 导入生成的 Swagger 文档
	"takeout/internal/router"
)

func routerInit() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // 接口文档
	r.StaticFS("/static", http.Dir("static"))                            //本地存储

	allRouter := router.AllRouter

	// admin
	admin := r.Group("/admin")
	{
		allRouter.EmployeeRouter.InitApiRouter(admin) // 注册员工路由
		allRouter.CategoryRouter.InitApiRouter(admin) // 注册菜品类别路由
		allRouter.DishRouter.InitApiRouter(admin)     // 注册菜品路由
		allRouter.CommonRouter.InitApiRouter(admin)   // 注册文件上传路由
		allRouter.SetMealRouter.InitApiRouter(admin)  // 注册套餐路由
		//allRouter.ShopRouter.InitApiRouter(admin) // 注册商店路由

		//allRouter.OrderRouter.InitApiRouter(admin)// 注册订单路由
		//allRouter.reportRouter.InitApiRouter(admin)// 注册报表路由
		//allRouter.workSpaceRouter.InitApiRouter(admin)// 注册工作台路由
	}
	user := r.Group("/user")
	{
		allRouter.WxUserRouter.InitApiRouter(user)
		allRouter.UserCategory.InitApiRouter(user)
		allRouter.UserShop.InitApiRouter(user)
		allRouter.UserDish.InitApiRouter(user)
		allRouter.UserSetmeal.InitApiRouter(user)
		allRouter.UserAddressBook.InitApiRouter(user)
		allRouter.UserShoppingCart.InitApiRouter(user)
		allRouter.UserOrder.InitApiRouter(user)
	}
	return r
}
