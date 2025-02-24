package router

import (
	"takeout/internal/router/admin"
	"takeout/internal/router/user"
	"takeout/internal/router/websocket"
)

type RouterGroup struct {
	admin.EmployeeRouter
	admin.CategoryRouter
	admin.DishRouter
	admin.CommonRouter
	admin.SetMealRouter
	admin.ShopRouter
	admin.OrderRouter
	admin.ReportRouter
	admin.WorkSpaceRouter
	websocket.Server
	UserWxUserRouter user.WxUserRouter
	UserShop         user.ShopRouter
	UserCategory     user.CategoryRouter
	UserDish         user.DishRouter
	UserSetmeal      user.SetmealRouter
	UserShoppingCart user.ShoppingCartRouter
	UserAddressBook  user.AddressBookRouter
	UserOrder        user.OrderRouter
}

var AllRouter = new(RouterGroup)
