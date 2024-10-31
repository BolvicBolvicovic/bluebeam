package analyzer

import (
	"net/url"
	"log"
	"sync"
	"github.com/gocolly/colly"
	"fmt"
)

type PageData struct {
	URL		string `json:"pageurl"`
	HTMLContent	string `json:"pagecontent"`
}

type Website struct {
	Pages		[]PageData `json:"pages"`
	mutex		sync.Mutex
}

func crawlWebsite(rootURL string, crawledWebsites *WebsiteGroup, wgUrls *sync.WaitGroup) {
	defer wgUrls.Done()

	var currentWebsite Website
	
	u, err := url.Parse(rootURL)
	if err != nil {
		log.Println("Invalid URL:", err)
		return
	}
	domain := u.Host
	
	c := colly.NewCollector(
		colly.AllowedDomains(domain),
		colly.Async(true),
		colly.MaxDepth(2),
		colly.MaxBodySize(20000),
	)

	visited := make(map[string]struct{})
	var visitedMutex sync.Mutex

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error during crawl:", err, string(r.Body))
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		u, err := e.Request.URL.Parse(link)
		if err != nil {
			log.Println("Invalid URL:", err)
			return
		}
		finalURL := u.String()

		visitedMutex.Lock()
		if _, found := visited[finalURL]; !found {
			visited[finalURL] = struct{}{}
			e.Request.Visit(finalURL)
		}
		visitedMutex.Unlock()
	})

	c.OnResponse(func(r *colly.Response) {
		// Process each page's content here
		if r.StatusCode >= 400 {
			return
		} 
		page := PageData{
			URL:         r.Request.URL.String(),
			HTMLContent: string(r.Body),
		}
		
		currentWebsite.mutex.Lock()
		currentWebsite.Pages = append(currentWebsite.Pages, page)
		currentWebsite.mutex.Unlock()
		fmt.Println("Visited", r.Request.URL)
	})

	fmt.Println("Starting crawl at:", rootURL)
	if err := c.Visit(rootURL); err != nil {
		log.Println("Error on start of crawl:", err)
	}

	c.Wait()
	crawledWebsites.mutex.Lock()
	crawledWebsites.Websites = append(crawledWebsites.Websites, currentWebsite)
	crawledWebsites.mutex.Unlock()
}
