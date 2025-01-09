package request

import (
	"encoding/json"
	"fmt"
	"strconv"
	"takeout/internal/model"
)

type SetMealDTO struct {
	Id           uint64              `json:"id"`            // 主键id
	CategoryId   uint64              `json:"categoryId"`    // 分类id
	Name         string              `json:"name"`          // 套餐名称
	Price        Price               `json:"price"`         // 套餐单价 前端存在bug，有时候发送string类型有时候是number类型
	Status       int                 `json:"status"`        // 套餐状态
	Description  string              `json:"description"`   // 套餐描述
	Image        string              `json:"image"`         // 套餐图片
	SetMealDishs []model.SetMealDish `json:"setmealDishes"` // 套餐菜品关系
}

type Price float64

// 实现 UnmarshalJSON 方法来处理 price 字段
func (p *Price) UnmarshalJSON(data []byte) error {
	// 尝试将数据解析为 float64
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*p = Price(f)
		return nil
	}

	// 如果是 string 类型，尝试将其解析为 float64
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		// 尝试将 string 转换为 float64
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("invalid price value: %v", s)
		}
		*p = Price(val)
		return nil
	}

	return fmt.Errorf("price should be a float64 or a string")
}

type SetMealPageQueryDTO struct {
	Page       int    `form:"page"`       // 分页查询的页数
	PageSize   int    `form:"pageSize"`   // 分页查询的页容量
	Name       string `form:"name"`       // 分页查询的name
	CategoryId uint64 `form:"categoryId"` // 分类ID:
	Status     string `form:"status"`     // 套餐起售状态
}
