package api

import (
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/bluebeam/templates/components"
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

	c.HTML(http.StatusOK, "settings.tmpl", gin.H{
		"username": username,
		"Navbar": components.NewNavbar(true),
	})
}

func MainPage(c *gin.Context) {
	isLoggedIn := true
	username, err := c.Cookie("bluebeam_username")
	if err != nil {
		isLoggedIn = false
	}
	session_key, err := c.Cookie("bluebeam_session_key")
	if err != nil {
		isLoggedIn = false
	}
	if username == "" || session_key == "" {
		isLoggedIn = false
	}
	if isLoggedIn {
		if !validUser(c, username, session_key) {
			return
		}
	}

	c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
		"Navbar": components.NewNavbar(isLoggedIn),
	})
}

func LoginPage(c *gin.Context) {
	isLoggedIn := true
	username, err := c.Cookie("bluebeam_username")
	if err != nil {
		isLoggedIn = false
	}
	session_key, err := c.Cookie("bluebeam_session_key")
	if err != nil {
		isLoggedIn = false
	}
	if username == "" || session_key == "" {
		isLoggedIn = false
	}
	if isLoggedIn {
		if !validUser(c, username, session_key) {
			return
		}
		c.Redirect(http.StatusOK, "/")
		return
	}
	c.HTML(http.StatusOK, "login_page.tmpl", gin.H{
		"Navbar": components.NewNavbar(false),
		"UsernameInput": components.Input{
			ID: "username",
		},
		"PasswordInput": components.Input{
			ID: "password",
		},
		"SubmitButton": components.Button{
			ID: "submitButton",
			Text: "login",
			IsSubmit: true,
		},
	})
}

func Logout(c *gin.Context) {
	username, err := c.Cookie("bluebeam_username")
	if err != nil {
		return
	}
	clearSessionKey(username)
	c.SetCookie("bluebeam_username", "", -1, "/", "localhost", true, true)
	c.SetCookie("bluebeam_session_key", "", -1, "/", "localhost", true, true)
	c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
		"Navbar": components.NewNavbar(false),
	})
}
