package request

type CategoryDTO struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Sort string `json:"sort"`
	Type string `json:"type"`
}

type CategoryPageQueryDTO struct {
	Name     string `form:"name" json:"name,omitempty"` // 分页查询的name
	Page     int    `form:"page" json:"page"`           // 分页查询的页数
	PageSize int    `form:"pageSize" json:"page_size"`  // 分页查询的页容量
	Cate     int    `form:"type" json:"cate,omitempty"` // 分类类型：1为菜品分类，2为套餐分类
}
