package api_test

import (
	"net/http"
	"net/http/httptest"
	"database/sql"
	"testing"
	"time"

	"github.com/BolvicBolvicovic/bluebeam/api"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.LoadHTMLGlob("templates/*/**")

	router.StaticFile("/favicon.ico", "./images/favicon.ico")
	router.StaticFile("/logo.png", "./images/logo.png")
	router.StaticFile("/example.png", "./images/example.png")
	router.StaticFile("/example2.png", "./images/example2.png")
	router.StaticFile("/main_page.md", "./templates/pages/main_page.md")
	router.StaticFile("/api_page.md", "./templates/pages/api_page.md")
	router.StaticFile("/why_bluebeam.md", "./templates/pages/why_bluebeam.md")
	router.StaticFile("/dashboard.js", "./templates/pages/dashboard.js")

	router.GET("/pong", api.Pong)
	router.GET("/dashboard", api.Dashboard)
	router.GET("/input_files", api.InputFiles)
	router.GET("/", api.MainPage)
	router.GET("/why_bluebeam", api.WhyBluebeam)
	router.GET("/api_page", api.ApiPage)
	router.GET("/login_page", api.LoginPage)
	router.GET("/logout", api.Logout)
	router.GET("/urls_output", api.UrlsOutput)

	return router
}

func setupMockDB(t *testing.T) sqlmock.Sqlmock {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	database.Db = db
	return mock
}

func registerAndLogin(username string, mock sqlmock.Sqlmock) {
	now := time.Now().Format(time.UnixDate)
	hash := "$2a$10$ipdWKjyptnKiOpLtb8w0gegE3bVGcIRr7lMzLD1nNVCqg5gTiByeq"
	mock.ExpectQuery("SELECT session_key, creation_key_time FROM users WHERE username = ?;").WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"session_key", "creation_key_time"}).
		AddRow(hash, now))
}

func TestPong(t *testing.T) {
	router := setupRouter()
	mock := setupMockDB(t)
	defer database.Db.Close()

	registerAndLogin("testUser", mock)
	req, _ := http.NewRequest("GET", "/pong", nil)
	req.AddCookie(&http.Cookie{
		Name: "bluebeam_username",
		Value: "testUser",
		MaxAge: 0,
		Secure: true,
		HttpOnly: true,
	})
	req.AddCookie(&http.Cookie{
		Name: "bluebeam_session_key",
		Value: "testPass",
		MaxAge: 0,
		Secure: true,
		HttpOnly: true,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "pong"}`, w.Body.String())
}

func TestDashboard(t *testing.T) {
	router := setupRouter()
	mock := setupMockDB(t)
	defer database.Db.Close()

	registerAndLogin("testUser", mock)
	req, _ := http.NewRequest("GET", "/dashboard", nil)
	req.AddCookie(&http.Cookie{
		Name: "bluebeam_username",
		Value: "testUser",
		MaxAge: 0,
		Secure: true,
		HttpOnly: true,
	})
	req.AddCookie(&http.Cookie{
		Name: "bluebeam_session_key",
		Value: "testPass",
		MaxAge: 0,
		Secure: true,
		HttpOnly: true,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "dashboard")
	assert.Contains(t, w.Body.String(), "popupOutput")
	assert.Contains(t, w.Body.String(), "let&#39;s go!") // Code for '
}

func TestInputFiles(t *testing.T) {
	router := setupRouter()
	mock := setupMockDB(t)
	defer database.Db.Close()

	registerAndLogin("testUser", mock)

	// Mock query for output_files_ids
	mock.ExpectQuery("SELECT criterias_files, current_file_index FROM users WHERE username = ?;").
		WithArgs("testUser").
		WillReturnError(sql.ErrNoRows)

	req, _ := http.NewRequest("GET", "/input_files", nil)
	req.AddCookie(&http.Cookie{
		Name: "bluebeam_username",
		Value: "testUser",
		MaxAge: 0,
		Secure: true,
		HttpOnly: true,
	})
	req.AddCookie(&http.Cookie{
		Name: "bluebeam_session_key",
		Value: "testPass",
		MaxAge: 0,
		Secure: true,
		HttpOnly: true,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"files": [], "index": -1}`, w.Body.String())
}

func TestLoginPage_NotLoggedIn(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/login_page", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "login")
}

func TestLogout(t *testing.T) {
	router := setupRouter()
	mock := setupMockDB(t)
	defer database.Db.Close()

	registerAndLogin("testUser", mock)

	// Mock clearSessionKey query
	mock.ExpectExec("UPDATE users SET session_key = NULL WHERE username = ?").
		WithArgs("testUser").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, _ := http.NewRequest("GET", "/logout", nil)
	req.AddCookie(&http.Cookie{
		Name: "bluebeam_username",
		Value: "testUser",
		MaxAge: 0,
		Secure: true,
		HttpOnly: true,
	})
	req.AddCookie(&http.Cookie{
		Name: "bluebeam_session_key",
		Value: "testPass",
		MaxAge: 0,
		Secure: true,
		HttpOnly: true,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

