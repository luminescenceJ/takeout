package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	OpenId     string    `json:"openid" gorm:"column:openid;unique;not null"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Sex        string    `json:"sex"`
	IdNumber   string    `json:"idNumber"`
	Avatar     string    `json:"avatar"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;type:datetime;not null"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreateTime = time.Now()
	return nil
}
