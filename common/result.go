package common

import (
	"gorm.io/gorm"
	"takeout/common/enum"
)

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type PageResult struct {
	Total   int64       `json:"total"`   // 总记录数
	Records interface{} `json:"records"` // 当前页数据集合
}

// PageVerify 分页查询 过滤器
func PageVerify(page *int, pageSize *int) {
	// 过滤 当前页、单页数量
	if *page < 1 {
		*page = 1
	}
	switch {
	case *pageSize > 100:
		*pageSize = enum.MaxPageSize
	case *pageSize <= 0:
		*pageSize = enum.MinPageSize
	}
}

func (p *PageResult) Paginate(page *int, pageSize *int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		PageVerify(page, pageSize)                                 // 校验参数
		return db.Offset(*pageSize * (*page - 1)).Limit(*pageSize) // 返回链式调用
	}
}
