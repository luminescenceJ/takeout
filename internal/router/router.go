package router

import (
	"takeout/internal/router/admin"
	"takeout/internal/router/user"
)

type RouterGroup struct {
	admin.EmployeeRouter
	admin.CategoryRouter
	admin.DishRouter
	admin.CommonRouter
	admin.SetMealRouter
	admin.ShopRouter
	user.WxUserRouter
	UserShop     user.ShopRouter
	UserCategory user.CategoryRouter
	UserDish     user.DishRouter
	UserSetmeal  user.SetmealRouter
}

var AllRouter = new(RouterGroup)
