package database

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestCreateTables(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
	    log.Fatal("Error creating mock db:", err)
	}
	defer db.Close()
	Db = db
	
	
	// Simulate successful database connection using the mock
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS decrypt_keys").WillReturnResult(sqlmock.NewResult(1, 1))
	
	Shutdown := createTables()
	
	if err := mock.ExpectationsWereMet(); err != nil {
	    t.Errorf("There were unmet expectations: %v", err)
	}
	
	assert.NotNil(t, Shutdown)
	
	Shutdown()
	
	assert.NoError(t, db.Close(), "expected DB to be closed successfully")
}
