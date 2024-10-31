package api

import (
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"log"
	"net/http"
)

func CurrentInputFile(c *gin.Context) {
	username, err := c.Cookie("bluebeam_username")
	var newIndex struct {
		Index	string `json:"newindex"`
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Need to log in again"})
		return
	}
	session_key, err := c.Cookie("bluebeam_session_key")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Need to log in again"})
		return
	}
	if err := c.ShouldBindJSON(&newIndex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if !validUser(c, username, session_key) {
		return
	}
	
	query := `
UPDATE
	users
SET
	current_file_index = ?
WHERE
	username = ?;
	`

	if _, err := database.Db.Exec(query, newIndex.Index, username); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating index"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "index updated successfuly"})
}

func UpdateEmail(c *gin.Context) {
	var newEmail struct {
		Email		string `json:"email"`
	}
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
	if err := c.ShouldBindJSON(&newEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if !validUser(c, username, session_key) {
		return
	}
	query := `
UPDATE
	users
SET
	email = ?
WHERE
	username = ?;
	`
	if _, err := database.Db.Exec(query, newEmail.Email, username); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email updated successfuly"})
}