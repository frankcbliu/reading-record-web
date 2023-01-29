package routers

import (
	"fmt"
	"net/http"
	"reading-record-web/models"

	"github.com/gin-gonic/gin"
)

// 创建书架
func CreateBookSelf(c *gin.Context) {
	name := c.PostForm("name")

	response := gin.H{"code": 0}

	if name != "" {
		user, _ := c.Get("user_model")
		models.CreateBookSelfModel(name, uint32(user.(models.User).ID))
	} else {
		fmt.Println("[routers.book_self] name is empty.")
		response["code"] = 1
	}

	fmt.Println("[routers.book_self] response: ", response)
	c.JSON(http.StatusOK, response)
}

// 获取书架列表
func GetBookSelf(c *gin.Context) {
	data := make(map[string]uint32)
	response := gin.H{"code": 0, "data": data}
	user, _ := c.Get("user_model")
	models.GetBookSelfModel(uint32(user.(models.User).ID), &data)

	fmt.Println("[routers.book_self] response: ", response)
	c.JSON(http.StatusOK, response)
}
