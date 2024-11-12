package api_test

import (
//	"net/http"
	"net/http/httptest"
	"testing"
//	"bytes"
//	"encoding/json"

	"github.com/BolvicBolvicovic/bluebeam/api"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

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
