package api_test

import (
	"net/http"
	"net/http/httptest"
	"database/sql"
	"testing"
	"bytes"
	"encoding/json"

	"github.com/BolvicBolvicovic/bluebeam/api"
	"github.com/BolvicBolvicovic/bluebeam/criterias"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

func setAESKey(mock sqlmock.Sqlmock) {
	query := `
SELECT
	first_key
FROM
	decrypt_keys
WHERE
	id = 0;
	`
	mock.ExpectQuery(query).
		WillReturnError(sql.ErrNoRows)
	query = `
INSERT INTO
	decrypt_keys
	(first_key)
VALUES
	(?);
	`
	mock.ExpectExec(query).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	criterias.SetKey()
}

func TestClearSessionKey(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	mock := setupMockDB(t)
	defer database.Db.Close()
	query := `
UPDATE
	users
SET
	session_key = NULL,
	creation_key_time = NULL
WHERE
	username = ?;
	`

	mock.ExpectExec(query).
		WithArgs("testUser").
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err := api.ClearSessionKey("testUser", c)
	assert.Equal(t, nil, err)
}

func TestGetAndStoreCriteria(t *testing.T) {
	router := setupRouter()
	router.POST("/criterias", api.StoreCriterias)
	mock := setupMockDB(t)
	defer database.Db.Close()
	setAESKey(mock)
	registerAndLogin("testUser", mock)
// Get
	query := `
SELECT
	criterias_files,
	current_file_index
FROM
	users
WHERE
	username = ?;
	`
	mock.ExpectQuery(query).
		WithArgs("testUser").
		WillReturnRows(sqlmock.
			NewRows([]string{"criterias_files", "current_file_index"}).
			AddRow([]byte(""), -1,))

// Store	
	query = `
UPDATE
	users
SET
	criterias_files = ?,
	current_file_index = ?
WHERE
	username = ?;
	`

	mock.ExpectExec(query).
		WithArgs(sqlmock.AnyArg(), 0,"testUser").
		WillReturnResult(sqlmock.NewResult(1, 1))
	body := criterias.Criterias{
		Features: []criterias.Feature{},
		FileName: "testFile",
	}
	out, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/criterias", bytes.NewBuffer(out))
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
	assert.JSONEq(t, `{"message": "Criterias well recieved, Data processed!"}`, w.Body.String())
}

func TestLogin(t *testing.T) {
	router := setupRouter()
	router.POST("/login", api.Login)
	mock := setupMockDB(t)
	defer database.Db.Close()
	setAESKey(mock)
	hash := "$2a$10$ipdWKjyptnKiOpLtb8w0gegE3bVGcIRr7lMzLD1nNVCqg5gTiByeq"
	
	query := `
SELECT
	password
FROM
	users
WHERE
	username = ?;
	`
	mock.ExpectQuery(query).
		WithArgs("testUser").
		WillReturnRows(sqlmock.
			NewRows([]string{"password"}).
			AddRow(hash))
	
	query = `
UPDATE
	users
SET
	session_key = ?,
	creation_key_time = ?
WHERE
	username = ?;
	`
	mock.ExpectExec(query).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(),"testUser").
		WillReturnResult(sqlmock.NewResult(1, 1))

	body :=  struct {
		Username	string `json:"username"`
		Password	string `json:"password"`
	}{
		Username:	"testUser",
		Password:	"testPass",
	}
	out, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(out))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.JSONEq(t, `{"message": "connected!"}`, w.Body.String())
}

func TestRegister(t *testing.T) {
	router := setupRouter()
	router.POST("/register", api.ResgisterAccount)
	mock := setupMockDB(t)
	defer database.Db.Close()
	setAESKey(mock)
	
	query := `
SELECT
	username
FROM
	users
WHERE
	username = ?;
	`
	mock.ExpectQuery(query).
		WithArgs("testUser").
		WillReturnError(sql.ErrNoRows)
	
	query = `
INSERT INTO
	users
	(username, password, email, current_file_index)
VALUES
	(?, ?, ?, ?);
	`
	mock.ExpectExec(query).
		WithArgs("testUser", sqlmock.AnyArg(), "test@test.com", -1). // Have to use AnyArg since bcrypt will generate a different hash each time
		WillReturnResult(sqlmock.NewResult(1, 4))

	body :=  struct {
		Username	string `json:"username"`
		Password	string `json:"password"`
		Email		string `json:"email"`
	}{
		Username:	"testUser",
		Password:	"testPass",
		Email:		"test@test.com",
	}
	out, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(out))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.JSONEq(t, `{"message": "Account successfuly created"}`, w.Body.String())
}
