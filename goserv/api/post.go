package api

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/scraper/database"
	"database/sql"
	"crypto/rand"
	"encoding/base64"
	"log"
)


func Login(c *gin.Context) {
	var user struct {
		Username	string `json:"username"`
		Password	string `json:"password"`
	}	
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if database.Db == nil {
		log.Println("Database connection is not initialized.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}
	query := `
SELECT
	password
FROM
	users
WHERE
	username = ?
	`
	row := database.Db.QueryRow(query, user.Username)
	var pwd string
	if err := row.Scan(pwd); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}
	//TODO: Hash the password
	if pwd != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password"})
		return
	} else {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating key session"})
			return
		}
		strkey := base64.StdEncoding.EncodeToString(key)
		query = `
UPDATE
	users
SET
	session_key = ?
	creation_key_time = ?
WHERE
	username = ?
		`
		if _, err := database.Db.Exec(query, strkey, time.Now(), user.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating key session"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"session_key": strkey})
	}
}

func ResgisterAccount(c *gin.Context) {
	if database.Db == nil {
		log.Println("Database connection is not initialized.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}
	var newUser struct {
		Username	string `json:"username"`
		Password	string `json:"password"`
		//TODO: Add phone/email verification
		//TODO: On the frontend, double check the password and how strong it is
	}
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	query := `
SELECT
	username
FROM
	users
WHERE
	username = ?
	`
	row := database.Db.QueryRow(query, newUser.Username)
	var test string
	if err := row.Scan(test); err == sql.ErrNoRows {
		query = `
UPDATE
	users
SET
	username = ?
	password = ?
		`
		c.JSON(http.StatusAccepted, gin.H{"message": "Account successfuly created"})
	} else {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username already taken"})
	}
}
