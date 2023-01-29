package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"code": 400,
				"error": "无效的session"})
		return
	}
	session.Delete(username)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"code": 500,
				"error": "session存储错误"})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{"code": 200,
			"message": "Successfully logged out"})
}
