package api

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/scraper/database"
	"database/sql"
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"encoding/base64"
	"log"
)

func Analyse(c *gin.Context) {
	var scrapedData struct {
		Username	string `json:username`
		SessionKey	string `json:"sessionkey"`
		Links		[]string `json:"links"`
		Buttons		[]struct {
			Text	string `json:"text"`
			OnClick string `json:"onclick"`
		} `json:"buttons"`
		PageHtml	string `json:"pageHtml"`
	}
	if err := c.ShouldBindJSON(&scrapedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	query := `
SELECT
	session_key,
	creation_key_time,
FROM
	users
WHERE
	username = ?;
	`
	row := database.Db.QueryRow(query, scrapedData.Username)
	var sk string
	if err := row.Scan(&sk); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/key"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(sk), []byte(scrapedData.SessionKey)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/key"})
		return
	}

}

func Login(c *gin.Context) {
	var user struct {
		Username	string `json:"username"`
		Password	string `json:"password"`
	}	
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	query := `
SELECT
	password
FROM
	users
WHERE
	username = ?;
	`
	row := database.Db.QueryRow(query, user.Username)
	var pwd string
	if err := row.Scan(&pwd); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(user.Password)); err != nil {
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
		hash, err := bcrypt.GenerateFromPassword([]byte(strkey), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Hash error"})
			return
		}
		log.Println(strkey, "len:", len(strkey))
		query = `
UPDATE
	users
SET
	session_key = ?,
	creation_key_time = ?
WHERE
	username = ?;
		`
		if _, err := database.Db.Exec(query, hash, time.Now().Format(time.UnixDate), user.Username); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating key session"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"session_key": strkey})
	}
}

func ResgisterAccount(c *gin.Context) {
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
	if len(newUser.Password) > 20 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password too long, max 19 characters"})
		return
	}
	query := `
SELECT
	username
FROM
	users
WHERE
	username = ?;
	`
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hash error"})
		return
	}
	row := database.Db.QueryRow(query, newUser.Username)
	var test string
	if err := row.Scan(&test); err == sql.ErrNoRows {
		query = `
INSERT INTO
	users
	(username, password)
VALUES
	(?, ?);
		`
		if _, err := database.Db.Exec(query, newUser.Username, hash); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "Account successfuly created"})
	} else {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username already taken"})
	}
}
