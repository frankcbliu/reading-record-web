package models

import (
	"gorm.io/gorm"
)

// 书架
type BookSelf struct {
	gorm.Model
	Name   string // 书架名称
	UserId uint32
}

func CreateBookSelfModel(name string) {
	db := InitModel()
	db.Create(&BookSelf{Name: name})
}

// func GetBookSelfModel() map[string]uint32 {
// }

// todo:
// 利用中间件解决账号登录后的 token 解析问题
