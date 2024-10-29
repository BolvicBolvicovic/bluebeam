package api

import (
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/bluebeam/criterias"
	"github.com/BolvicBolvicovic/bluebeam/templates/components"
	"net/http"
)

//TODO: add UpdateEmail here

func CurrentInputFile(c *gin.Context) {
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
	
	query := `
UPDATE
	users
SET
	current_file_index = ?
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
