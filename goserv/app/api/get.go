package api

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/gin-gonic/gin"
	"net/http"
	"context"
)

var oauthConfig = &oauth2.Config{
	RedirectURL:  "https://localhost/selectGoogleFile",
	ClientID:     "726518157620-8s2194lb2ka65vfga9loee2sookpjfda.apps.googleusercontent.com",
	ClientSecret: "",
	Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets.readonly", "https://www.googleapis.com/auth/drive"},
	Endpoint:     google.Endpoint,
}

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

func InitOAuth(c *gin.Context) {
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, authURL)
}

func SelectGoogleFile(c *gin.Context) {
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
	if !validUser(c, username,session_key) {
		return
	}
	code := c.Query("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "selectGoogleFile.tmpl", gin.H {
		"token": token,
	})
}
