package models

import (
	"fmt"

	"gorm.io/gorm"
)

// 书架
type BookSelf struct {
	gorm.Model
	Name   string // 书架名称
	UserId uint32
}

func CreateBookSelfModel(name string, user_id uint32) bool {
	db := InitModel()
	if QueryBookSelfModel(name, user_id) {
		fmt.Println("[models.bookself] No need to recreate.")
		return false
	}
	db.Create(&BookSelf{Name: name, UserId: user_id})
	fmt.Println("[models.bookself]", name, "user:", user_id, "create.")
	return true
}

// 查询书架是否存在
func QueryBookSelfModel(name string, user_id uint32) bool {
	var book_self BookSelf
	db := InitModel()
	db.Where("name = ?", name).First(&book_self, "user_id = ?", user_id)
	if book_self.ID > 0 {
		fmt.Println("[models.bookself]", name, "user:", user_id, "exist")
		return true
	}
	return false
}

func GetBookSelfModel(user_id uint32, data *map[string]uint32) {
	var bookselves []BookSelf
	db := InitModel()

	db.Where("user_id = ?", user_id).Find(&bookselves)

	fmt.Println(len(bookselves))
	for _, v := range bookselves {
		(*data)[v.Name] = uint32(v.ID)
	}
}
