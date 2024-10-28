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

func Dashboard(c *gin.Context) {
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

	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
		"username": username,
		"Navbar": components.NewNavbar(true),
		"UrlsSubmitButton": components.Button{
			Text: "let's go!",
			IsSubmit: true,
		},
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

func WhyBluebeam(c *gin.Context) {
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

	c.HTML(http.StatusOK, "why_bluebeam.tmpl", gin.H{
		"Navbar": components.NewNavbar(isLoggedIn),
	})
}

func ApiPage(c *gin.Context) {
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

	c.HTML(http.StatusOK, "api_page.tmpl", gin.H{
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
		"UsernameInputLogin": components.Input{
			ID: "usernameLogin",
			Placeholder: "Username",
		},
		"PasswordInputLogin": components.Input{
			ID: "passwordLogin",
			Type: "password",
			Placeholder: "Password",
		},
		"SubmitButtonLogin": components.Button{
			ID: "submitButtonLogin",
			Text: "let's go!",
			IsSubmit: true,
		},
		"UsernameInputRegister": components.Input{
			ID: "usernameRegister",
			Placeholder: "Username",
		},
		"PasswordInputRegister": components.Input{
			ID: "passwordRegister",
			Type: "password",
			Placeholder: "Password",
		},
		"PasswordInputTester": components.Input{
			ID: "passwordTest",
			Type: "password",
			Placeholder: "re-type Password",
		},
		"EmailInputRegister": components.Input{
			ID: "email",
			Placeholder: "Email compatible with Google",
		},
		"SubmitButtonRegister": components.Button{
			ID: "submitButtonRegister",
			Text: "let's go!",
			IsSubmit: true,
		},
		"Switch": components.Button{
			ID: "switchButton",
			Text: "register",
			IsSubmit: false,
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
