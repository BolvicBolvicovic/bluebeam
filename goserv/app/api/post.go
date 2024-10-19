package api

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"github.com/BolvicBolvicovic/bluebeam/analyzer"
	"github.com/BolvicBolvicovic/bluebeam/criterias"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/api/drive/v3"
	"fmt"
	"database/sql"
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"encoding/base64"
	"encoding/json"
	"log"
)

func validUser(c *gin.Context, username string, session_key string) bool {
	query := `
SELECT
	session_key,
	creation_key_time
FROM
	users
WHERE
	username = ?;
	`
	row := database.Db.QueryRow(query, username)
	var sk, ckt sql.NullString
	if err := row.Scan(&sk, &ckt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/key"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return false
	}
	if !sk.Valid || !ckt.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/key"})
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(sk.String), []byte(session_key)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/key"})
		return false
	}
	creation_key_time, err := time.Parse(time.UnixDate, ckt.String)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return false
	}
	y0, m0, d0 := creation_key_time.Date()
	y1, m1, d1 := time.Now().Date()
	//TODO: Handle the session key hourly or with the ping function
	if d0 != d1 || m0 != m1 || y0 != y1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Outdated key"})
		return false
	}
	return true
}

func StoreCriterias(c *gin.Context) {
	crits := criterias.Criterias{}
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
	if err := c.ShouldBindJSON(&crits); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if !validUser(c, username, session_key) {
		return
	}
	criterias.Store(c, crits, username)
}

func Analyze(c *gin.Context) {
	var scrapedData analyzer.ScrapedDefault
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
	if err := c.ShouldBindJSON(&scrapedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if !validUser(c, username, session_key) {
		return
	}
	analyzer.Analyzer(c, scrapedData, username)
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

func createAndFillSpreadsheet(srv *sheets.Service, d *drive.Service, data json.RawMessage, email string) (string, error) {
	spreadsheet := &sheets.Spreadsheet{
	    Properties: &sheets.SpreadsheetProperties{
	        Title: "New dataset",
	    },
	}
	
	sheet, err := srv.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
	    return "", fmt.Errorf("unable to create spreadsheet: %v", err)
	}
	
	spreadsheetID := sheet.SpreadsheetId
	
	// Parse the JSON data
	var records []map[string]interface{}
	if err := json.Unmarshal(data, &records); err != nil {
	    return "", fmt.Errorf("error parsing JSON: %v", err)
	}
	
	// Check if records are available
	if len(records) == 0 {
	    return "", fmt.Errorf("no data to fill")
	}
	
	// Extract headers from the first record
	headers := []string{}
	for key := range records[0] {
	    headers = append(headers, key)
	}
	
	// Prepare the 2D array for Google Sheets
	values := [][]interface{}{}
	
	// Add headers as the first row
	headerRow := make([]interface{}, len(headers))
	for i, header := range headers {
	    headerRow[i] = header
	}
	values = append(values, headerRow)
	
	// Convert each JSON object into a row
	for _, record := range records {
	    row := make([]interface{}, len(headers))
	    for i, header := range headers {
	        row[i] = record[header]
	    }
	    values = append(values, row)
	}
	
	vr := &sheets.ValueRange{
	    Values: values,
	}
	
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, "Sheet1!A1", vr).ValueInputOption("RAW").Do()
	if err != nil {
	    return "", fmt.Errorf("unable to update spreadsheet values: %v", err)
	}

	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer", // Can be "reader" or "writer"
		EmailAddress: email,
	}
	_, err = d.Permissions.Create(spreadsheetID, permission).SendNotificationEmail(true).Do()
	if err != nil {
	        return "", fmt.Errorf("unable to share spreadsheet: %v", err)
	}

	return fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", spreadsheetID), nil
}

func OutputGoogleSpreadsheet(c *gin.Context) {
	var output struct {
		Data		json.RawMessage `json:"data"`
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
	if err := c.ShouldBindJSON(&output); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if !validUser(c, username, session_key) {
		return
	}
	query := `
SELECT
	email
FROM
	users
WHERE
	username = ?
	`
	row := database.Db.QueryRow(query, username)
	var email sql.NullString
	if err := row.Scan(&email); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}
	if !email.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}
	service, exists := c.Get("sheetsService")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sheets service not found"})
		return
	}
	sheetsService, ok := service.(*sheets.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Sheets service instance"})
		return
	}
	service, exists = c.Get("driveService")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Drive service not found"})
		return
	}
	driveService, ok := service.(*drive.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid drive service instance"})
		return
	}

	url, err := createAndFillSpreadsheet(sheetsService, driveService, output.Data, email.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"spreadsheetUrl": url})
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
	var pwd sql.NullString
	if err := row.Scan(&pwd); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}
	if !pwd.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(pwd.String), []byte(user.Password)); err != nil {
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
		c.SetCookie("bluebeam_username", user.Username, 3600, "/", "localhost", true, true)
		c.SetCookie("bluebeam_session_key", strkey, 3600, "/", "localhost", true, true)
		c.JSON(http.StatusAccepted, gin.H{"message": "connected!"})
	}
}

func ResgisterAccount(c *gin.Context) {
	var newUser struct {
		Username	string `json:"username"`
		Password	string `json:"password"`
		Email		string `json:"email"`
		//TODO: Add email verification and hash usernam/email
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
	(username, password, email)
VALUES
	(?, ?, ?);
		`
		if _, err := database.Db.Exec(query, newUser.Username, hash, newUser.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "Account successfuly created"})
	} else {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username already taken"})
	}
}
