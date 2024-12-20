package router

import "takeout/internal/router/admin"

type RouterGroup struct {
	admin.EmployeeRouter
}

var AllRouter = new(RouterGroup)
