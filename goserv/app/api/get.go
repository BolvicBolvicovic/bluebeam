package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H { "message": "pong", })
}

func Settings(c *gin.Context) {
	var user struct {
		Username	string `form:"username" binding:"required"`
		SessionKey	string `form:"sessionkey" binding:"required"`
	}

	if err := c.ShouldBindQuery(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validUser(c, user.Username, user.SessionKey) {
		return
	}

	c.HTML(http.StatusOK, "settings.tmpl", gin.H {
		"username": user.Username,
		"sessionkey": user.SessionKey,
	})
}
