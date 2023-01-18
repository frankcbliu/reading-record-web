package models

import (
	"crypto/md5"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string
	Password string
}

// 创建用户
func CreateUserModel(user_name string, password string) bool {
	db := InitModel()
	var user User
	if !QueryUserModel(&user, user_name) {
		return false
	}
	db.Create(&User{UserName: user_name, Password: MD5(password)})
	return true
}

// 用户登录
func LoginUserModel(user_name string, password string) bool {
	var user User
	if !QueryUserModel(&user, user_name) {
		return false
	}
	if user.Password != MD5(password) {
		fmt.Println(user_name + " password error!")
		return false
	}
	return true
}

// 用户查询
func QueryUserModel(user *User, user_name string) bool {
	db := InitModel()
	db.First(user, "user_name = ?", user_name)
	if user.ID > 0 {
		fmt.Println(user.UserName, "exist")
		return true
	}
	fmt.Println(user.UserName, "not exist")
	return false
}

// md5 hash
func MD5(input string) string {
	hash := md5.Sum([]byte("reading-record-web" + input))
	return fmt.Sprintf("%x", hash)
}
