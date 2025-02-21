package model

import "time"

// ShoppingCart 购物车数据模型
type ShoppingCart struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	UserId     int       `json:"userId"`
	DishId     int       `json:"dishId"`
	SetmealId  int       `json:"setmealId"`
	DishFlavor string    `json:"dishFlavor"`
	Number     int       `json:"number"`
	Amount     float64   `json:"amount"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (ShoppingCart) TableName() string {
	return "shopping_cart"
}
