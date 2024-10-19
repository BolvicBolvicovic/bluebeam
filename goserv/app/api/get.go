package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


func Pong(c *gin.Context) {
	username, err := c.Cookie("bluebeam_username")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Need to log in again"})
		return
	}
	session_key, err := c.Cookie("bluebeam_session_key")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Need to log in again"})
		return
	}

	if !validUser(c, username, session_key) {
		return
	}
	c.JSON(http.StatusOK, gin.H { "message": "pong", })
}

func Settings(c *gin.Context) {
	username, err := c.Cookie("bluebeam_username")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Need to log in again"})
		return
	}
	session_key, err := c.Cookie("bluebeam_session_key")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Need to log in again"})
		return
	}

	if !validUser(c, username, session_key) {
		return
	}

	c.HTML(http.StatusOK, "settings.tmpl", gin.H{})
}
