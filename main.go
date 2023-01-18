package main

import (
	"reading-record-web/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 用户相关
	user := r.Group("/user")
	{
		user.POST("/register", routers.RegisterHandler)
		user.POST("/login", routers.LoginHandler)
	}

	// // 读取内容
	// var product Product
	// db.First(&product, 1)                 // find product with integer primary key
	// db.First(&product, "code = ?", "D42") // find product with code D42

	// // 更新操作： 更新单个字段
	// db.Model(&product).Update("Price", 2000)

	// // 更新操作： 更新多个字段
	// db.Model(&product).Updates(Product{Price: 2000, Code: "F42"}) // non-zero fields
	// db.Model(&product).Updates(map[string]interface{}{"Price": 2000, "Code": "F42"})

	// // 删除操作：
	// db.Delete(&product, 1)

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
