package router

import "takeout/internal/router/admin"

type RouterGroup struct {
	admin.EmployeeRouter
	admin.CategoryRouter
	admin.DishRouter
}

var AllRouter = new(RouterGroup)
