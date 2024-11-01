package api

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/bluebeam/templates/components"
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

func clearSessionKey(username string, c *gin.Context) error {
	query := `
UPDATE
	users
SET
	session_key = NULL,
	creation_key_time = NULL
WHERE
	username = ?;
	`
	c.SetCookie("bluebeam_username", "", -1, "/", "localhost", true, true)
	c.SetCookie("bluebeam_session_key", "", -1, "/", "localhost", true, true)
	_, err := database.Db.Exec(query, username)
	return err
}

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
			clearSessionKey(username, c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/key"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return false
	}
	if !sk.Valid || !ckt.Valid {
		clearSessionKey(username, c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/key"})
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(sk.String), []byte(session_key)); err != nil {
		clearSessionKey(username, c)
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
		clearSessionKey(username, c)
		c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
			"Navbar": components.NewNavbar(false),
		})
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if !validUser(c, username, session_key) {
		return
	}
	criterias.Store(c, crits, username)
}

func Urls(c *gin.Context) {
	var scrapedUrls analyzer.ScrapedUrls
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
	if err := c.ShouldBindJSON(&scrapedUrls); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if !validUser(c, username, session_key) {
		return
	}
	analyzer.HandleUrls(c, scrapedUrls, username)
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

type ResponseJSON struct {
	FeatureName    string `json:"feature_name"`
	IsPresent      bool   `json:"ispresent"`
	TextIfPresent  string `json:"textifpresent"`
	ThoughtProcess string `json:"thoughtprocess"`
}

func createAndFillSpreadsheet(srv *sheets.Service, d *drive.Service, data json.RawMessage, email string) (string, error) {
	// Create the spreadsheet
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: "New dataset",
		},
	}

	// Create a new spreadsheet
	sheet, err := srv.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create spreadsheet: %v", err)
	}

	spreadsheetID := sheet.SpreadsheetId
	// Parse the JSON data into an array of responses
	var responses [][]ResponseJSON
	if err := json.Unmarshal(data, &responses); err != nil {
		return "", fmt.Errorf("error parsing JSON: %v", err)
	}

	// Iterate over each response set to create a new sheet
	for i, responseSet := range responses {
		sheetTitle := fmt.Sprintf("Response %d", i+1)
		_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{
					AddSheet: &sheets.AddSheetRequest{
						Properties: &sheets.SheetProperties{
							Title: sheetTitle,
						},
					},
				},
			},
		}).Do()
		if err != nil {
			return "", fmt.Errorf("unable to create sheet %s: %v", sheetTitle, err)
		}

		// Prepare headers and values for the new sheet
		headers := []string{"FeatureName", "IsPresent", "TextIfPresent", "ThoughtProcess"}
		values := [][]interface{}{
			{headers[0], headers[1], headers[2], headers[3]},
		}

		// Populate the rows from each ResponseJSON item in responseSet
		for _, record := range responseSet {
			row := []interface{}{record.FeatureName, record.IsPresent, record.TextIfPresent, record.ThoughtProcess}
			values = append(values, row)
		}

		// Write the values to the new sheet
		vr := &sheets.ValueRange{
			Values: values,
		}
		_, err = srv.Spreadsheets.Values.Update(spreadsheetID, fmt.Sprintf("%s!A1", sheetTitle), vr).ValueInputOption("RAW").Do()
		if err != nil {
			return "", fmt.Errorf("unable to update values in sheet %s: %v", sheetTitle, err)
		}
	}

	// Share the spreadsheet with the specified email
	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: email,
	}
	_, err = d.Permissions.Create(spreadsheetID, permission).SendNotificationEmail(true).Do()
	if err != nil {
		return "", fmt.Errorf("unable to share spreadsheet: %v", err)
	}

	return fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", spreadsheetID), nil
}

func storeUrlOutput(url string, username string) {
	urls, err := getUrls(username)
	if err != nil {
		log.Println("Error getting urls:", err)
		return
	}
	urls = append(urls, []byte(url)...)
	query := `
UPDATE
	users
SET
	output_files_ids = ?
WHERE
	username = ?;
	`
	if _, err := database.Db.Exec(query, []byte(urls), username); err != nil {
		log.Println(err)
	}
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
		fmt.Println(err)
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
	storeUrlOutput(url, username)
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
		c.SetCookie("bluebeam_username", user.Username, 86400, "/", "localhost", true, true)
		c.SetCookie("bluebeam_session_key", strkey, 86400, "/", "localhost", true, true)
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
	(username, password, email, current_file_index)
VALUES
	(?, ?, ?, ?);
		`
		if _, err := database.Db.Exec(query, newUser.Username, hash, newUser.Email, -1); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "Account successfuly created"})
	} else {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username already taken"})
	}
}
