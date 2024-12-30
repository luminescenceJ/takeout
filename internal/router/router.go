package router

import "takeout/internal/router/admin"

type RouterGroup struct {
	admin.EmployeeRouter
	admin.CategoryRouter
}

var AllRouter = new(RouterGroup)
