package api

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"context"
)

var (
	oauthConfig = &oauth2.Config{
		RedirectURL:  "https://localhost/selectGoogleFile",
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets.readonly"},
		Endpoint:     google.Endpoint,
	}
	tokenStore = map[string]*oauth2.Token{} // Temporary store; replace with secure storage for production
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

func InitOAuth(c *gin.Context) {

}

func SelectGoogleFile(c *gin.Context) {

}
