package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 初始化模型
func InitModel() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&User{})
	return db
}
