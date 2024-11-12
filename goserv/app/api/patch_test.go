package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"encoding/json"

	"github.com/BolvicBolvicovic/bluebeam/api"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCurrentInputFile(t *testing.T) {
	router := setupRouter()
	mock := setupMockDB(t)
	defer database.Db.Close()
	
	router.PATCH("/current_input_file", api.CurrentInputFile)

	registerAndLogin("testUser", mock)
	
	query := `
UPDATE
	users
SET
	current_file_index = ?
WHERE
	username = ?;
	`

	mock.ExpectExec(query).
		WithArgs("1", "testUser").
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := struct { Index	string `json:"newindex"` } { Index: "1", }
	out, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", "/current_input_file", bytes.NewBuffer(out))
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
	assert.JSONEq(t, `{"message": "index updated successfuly"}`, w.Body.String())
}

func TestUpdateEmail(t *testing.T) {
	router := setupRouter()
	mock := setupMockDB(t)
	defer database.Db.Close()
	
	router.PATCH("/email", api.UpdateEmail)

	registerAndLogin("testUser", mock)
	
	query := `
UPDATE
	users
SET
	email = ?
WHERE
	username = ?;
	`

	mock.ExpectExec(query).
		WithArgs("test@test.com", "testUser").
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := struct { Email	string `json:"email"` } { Email: "test@test.com", }
	out, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", "/email", bytes.NewBuffer(out))
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
	assert.JSONEq(t, `{"message": "Email updated successfuly"}`, w.Body.String())
}

func TestUpdateAPIKey(t *testing.T) {
	router := setupRouter()
	mock := setupMockDB(t)
	defer database.Db.Close()
	
	router.PATCH("/api_key", api.UpdateAPIKey)

	registerAndLogin("testUser", mock)
	
	query := `
UPDATE
	users
SET
	openai_api_key = ?
WHERE
	username = ?;
	`

	mock.ExpectExec(query).
		WithArgs("mocked_key", "testUser").
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := struct { 
		Type	string `json:"type"`
		APIKey	string `json:"apikey"`
	} { 
		Type: "openai",
		APIKey: "mocked_key",
	}
	out, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", "/api_key", bytes.NewBuffer(out))
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
	assert.JSONEq(t, `{"message": "API key updated successfuly"}`, w.Body.String())
}
