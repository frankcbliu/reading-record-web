package routers

import (
	"net/http"
	"reading-record-web/models"

	"github.com/gin-gonic/gin"
)

// 注册接口
func RegisterHandler(c *gin.Context) {
	user_name := c.Param("user_name")
	password := c.Param("password")
	response := gin.H{"code": 0}
	if !models.CreateUserModel(user_name, password) {
		response["code"] = 1
	}
	c.JSON(http.StatusOK, response)
}

// 登录接口
func LoginHandler(c *gin.Context) {
	user_name := c.Param("user_name")
	password := c.Param("password")
	response := gin.H{"code": 0}
	if !models.LoginUserModel(user_name, password) {
		response["code"] = 1
	}
	c.JSON(http.StatusOK, response)
}
