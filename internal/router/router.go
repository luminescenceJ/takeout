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
}

var AllRouter = new(RouterGroup)
