package routers

import (
	"fmt"
	"net/http"
	"reading-record-web/models"
	"reading-record-web/tools"

	"github.com/gin-gonic/gin"
)

// 注册接口
func RegisterHandler(c *gin.Context) {
	user_name := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println("[routers.user][Register] user_name:", user_name)
	response := gin.H{"code": 0}
	if !models.CreateUserModel(user_name, password) {
		response["code"] = 1
	}
	fmt.Println("[routers.user][Register] response: ", response)
	c.JSON(http.StatusOK, response)
}

// 登录接口
func LoginHandler(c *gin.Context) {
	user_name := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println("[routers.user][Login] user_name:", user_name)

	response := gin.H{"code": 0}
	if !models.LoginUserModel(user_name, password) {
		response["code"] = 1
	} else {
		response["token"] = tools.GenToken(user_name)
	}
	fmt.Println("[routers.user][Login] response: ", response)
	c.JSON(http.StatusOK, response)
}
