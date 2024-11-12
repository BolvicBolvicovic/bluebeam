package analyzer

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test crawlWebsite with a mock server
func TestCrawlWebsite(t *testing.T) {
	// Start a mock server to simulate a webpage
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<head><title>Test Page</title></head>
				<body>
					<a href="/page1">Link to Page 1</a>
					<a href="/page2">Link to Page 2</a>
					<img src="image.jpg" alt="Test Image"/>
					<h1>Page Header</h1>
					<button onclick="alert('Clicked!')">Test Button</button>
				</body>
			</html>
		`))
	}))
	defer mockServer.Close()

	// Initialize a WebsiteGroup to store crawled data
	var wg sync.WaitGroup
	var crawledWebsites WebsiteGroup

	// Call the crawlWebsite function with the mock server URL
	wg.Add(1)
	crawlWebsite(mockServer.URL, &crawledWebsites, &wg)
	wg.Wait()

	// Assertions to verify the crawling process
	assert.Len(t, crawledWebsites.Websites, 1)              // One website should be added
	assert.Greater(t, len(crawledWebsites.Websites[0]), 0)  // At least one page should be crawled

	// Check contents of the first crawled page
	page := crawledWebsites.Websites[0][0]
	assert.Equal(t, mockServer.URL, page.URL)
	assert.Contains(t, page.HTMLContent.BodyInnerText, "Link to Page 1")
	assert.Contains(t, page.HTMLContent.BodyInnerText, "Link to Page 2")
	assert.Contains(t, page.HTMLContent.BodyInnerText, "Test Button")
}

// Test parsePage directly with a sample HTML body
func TestParsePage(t *testing.T) {
	htmlBody := `
		<html>
			<head><title>Test Page</title></head>
			<body>
				<a href="https://example.com/page1">Link to Page 1</a>
				<a href="https://example.com/page2">Link to Page 2</a>
				<img src="image.jpg" alt="Test Image"/>
				<h1>Main Header</h1>
				<button onclick="alert('Clicked!')">Test Button</button>
			</body>
		</html>
	`

	pageData, err := parsePage("https://example.com", htmlBody)
	assert.NoError(t, err)

	// Verify parsed data
	assert.Equal(t, "https://example.com", pageData.URL)
	assert.Contains(t, pageData.HTMLContent.BodyInnerText, "Link to Page 1")
	assert.Contains(t, pageData.HTMLContent.BodyInnerText, "Link to Page 2")
	assert.Contains(t, pageData.HTMLContent.BodyInnerText, "Main Header")
	assert.Contains(t, pageData.HTMLContent.BodyInnerText, "Test Button")
	assert.Len(t, pageData.HTMLContent.Images, 1)
	assert.Equal(t, "image.jpg", pageData.HTMLContent.Images[0].Src)
	assert.Equal(t, "Test Image", pageData.HTMLContent.Images[0].Alt)
}
