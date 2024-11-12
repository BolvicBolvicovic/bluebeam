package api

import (
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/bluebeam/criterias"
	"github.com/BolvicBolvicovic/bluebeam/templates/components"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"database/sql"
	"log"
	"errors"
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
		"PopupOutput": components.NewPopupOutput("popupOutput"),
		"UrlsSubmitButton": components.Button{
			Text: "let's go!",
			IsSubmit: true,
		},
		"InputChoiceSubmitButton": components.Button{
			Text: "update input file",
			IsSubmit: true,
		},
	})
}

func InputFiles(c *gin.Context) {
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
	inputFiles, index_file, err := criterias.Get(c, username)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H {
		"files": inputFiles,
		"index": index_file,
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
	clearSessionKey(username, c)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func getUrls(username string) ([]byte, error) {
	query := `
SELECT
	output_files_ids	
FROM
	users
WHERE
	username = ?
	`
	row := database.Db.QueryRow(query, username)
	var urls sql.Null[[]byte]
	if err := row.Scan(&urls); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid username")
		} else {
			log.Println(err)
			return nil, errors.New("Internal error")
		}
	}
	if !urls.Valid {
		return make([]byte, 0), nil
	}
	return []byte(urls.V), nil
}

func UrlsOutput(c *gin.Context) {
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
	urlsOutput, err := getUrls(username)
	if err != nil {
		log.Println("UrlsOutput:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"urlsoutput": urlsOutput})
}
