package analyzer

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/BolvicBolvicovic/bluebeam/criterias"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAPIKey(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.Db = db

	t.Run("Success - Valid Key", func(t *testing.T) {
		expectedKey := "mocked_api_key"
		rows := sqlmock.NewRows([]string{"openai_api_key"}).AddRow(expectedKey)
		mock.ExpectQuery("SELECT openai_api_key FROM users WHERE username = ?").WithArgs("user123").WillReturnRows(rows)

		key, err := getAPIKey("OPENAI_API_KEY", "user123")
		assert.Nil(t, err)
		assert.Contains(t, key, expectedKey)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failure - No Key Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT openai_api_key FROM users WHERE username = ?").WithArgs("nonexistent_user").WillReturnError(sql.ErrNoRows)

		_, err := getAPIKey("OPENAI_API_KEY", "nonexistent_user")
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func mockCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestSendLLMQuestion(t *testing.T) {
	mockFeature := criterias.Feature{FeatureName: "Test Feature"}
	var jsonData = json.RawMessage(`{"data": "test"}`)
	response := &LLMResponse{}
	var wg sync.WaitGroup

	t.Run("Failure - Invalid AI Client Path", func(t *testing.T) {
		wg.Add(1)
		sendLLMQuestion(mockFeature, &jsonData, response, &wg, "nonexistent_ai", "OPENAI_API_KEY=mocked_key")
		wg.Wait()

		assert.Empty(t, response.Response)
	})

	t.Run("Success - Executes Command and Updates Response", func(t *testing.T) {
		wg.Add(1)
		sendLLMQuestion(mockFeature, &jsonData, response, &wg, "gpt-4o-mini", "OPENAI_API_KEY=mocked_key")
		wg.Wait()

		assert.NotEmpty(t, response.Response)
	})
}

func TestSanitizeCrawledWebsites(t *testing.T) {
	t.Run("Success - Returns Sanitized Data", func(t *testing.T) {
		mockGroup := WebsiteGroup{
			Websites: [][]PageData{
				{{URL: "Test Page"}},
			},
		}
		sanitizedData, err := sanitizeCrawledWebsites(mockGroup, "gpt-4o-mini", "OPENAI_API_KEY=mocked_key", "default")
		assert.Nil(t, err)
		assert.NotEmpty(t, sanitizedData.Websites)
	})

	t.Run("Failure - Invalid AI Client", func(t *testing.T) {
		mockGroup := WebsiteGroup{}
		_, err := sanitizeCrawledWebsites(mockGroup, "nonexistent_ai", "OPENAI_API_KEY=mocked_key", "default")
		assert.NotNil(t, err)
	})
}

func TestAnalyzer(t *testing.T) {
	router := gin.Default()
	router.POST("/analyze", func(c *gin.Context) {
		username := "testuser"
		var sd ScrapedDefault
		Analyzer(c, sd, username)
	})

	w := performRequest(router, "POST", "/analyze")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No API key for openai")
}

func TestHandleUrls(t *testing.T) {
	router := gin.Default()
	router.POST("/handle-urls", func(c *gin.Context) {
		username := "testuser"
		su := ScrapedUrls{Ai: "gpt-4o-mini", Sanitizer: "default"}
		HandleUrls(c, su, username)
	})

	w := performRequest(router, "POST", "/handle-urls")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No API key for the chosen AI")
}

// Helper function to perform HTTP requests
func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

