package analyzer

import (
	"net/url"
	"golang.org/x/net/html"
	"strings"
	"log"
	"sync"
	"github.com/gocolly/colly"
	"fmt"
)

type PageData struct {
	URL		string		`json:"pageurl"`
	HTMLContent	ScrapedDefault  `json:"pagecontent"`
}

type Website struct {
	Pages		[]PageData `json:"pages"`
	mutex		sync.Mutex
}

type ScrapedDefaultMutex struct {
	MutexLinks		sync.Mutex
	MutexButtons		sync.Mutex 
	MutexImages		sync.Mutex
	MutexFormInputs		sync.Mutex
	MutexMetaTags		sync.Mutex 
	MutexHeaders		sync.Mutex 
	MutexBodyInnerText	sync.Mutex
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
		colly.MaxBodySize(60000),
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
		defer fmt.Println("Visited", r.Request.URL)

		if r.StatusCode >= 400 {
			return
		} 
		page, err := parsePage(r.Request.URL.String(), string(r.Body))
		if err != nil {
			log.Println("At:", r.Request.URL, "error parsing page:", err)
			return
		}
		
		currentWebsite.mutex.Lock()
		currentWebsite.Pages = append(currentWebsite.Pages, page)
		currentWebsite.mutex.Unlock()
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

func processNode(n *html.Node, content *ScrapedDefault, mutexes *ScrapedDefaultMutex, wg *sync.WaitGroup) {
	defer wg.Done()
	if n.Type == html.ElementNode {
		switch n.Data {
		case "script", "style", "noscript":
			return
		case "a":
			var link _Link
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link.Href = attr.Val
				}
			}
			if n.FirstChild != nil {
				link.Text = n.FirstChild.Data
			}
			mutexes.MutexLinks.Lock()
			content.Links = append(content.Links, link)
			mutexes.MutexLinks.Unlock()

		case "button":
			var button _Button
			for _, attr := range n.Attr {
				if attr.Key == "onclick" {
					button.OnClick = attr.Val
				} else if attr.Key == "id" {
					button.ID = attr.Val
				} else if attr.Key == "class" {
					button.Classes = attr.Val
				}
			}
			if n.FirstChild != nil {
				button.Text = n.FirstChild.Data
			}
			mutexes.MutexButtons.Lock()
			content.Buttons = append(content.Buttons, button)
			mutexes.MutexButtons.Unlock()

		case "img":
			var image _Image
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					image.Src = attr.Val
				} else if attr.Key == "alt" {
					image.Alt = attr.Val
				} else if attr.Key == "class" {
					image.Classes = attr.Val
				}
			}
			mutexes.MutexImages.Lock()
			content.Images = append(content.Images, image)
			mutexes.MutexImages.Unlock()

		case "input":
			var input _FormInput
			for _, attr := range n.Attr {
				if attr.Key == "type" {
					input.Type = attr.Val
				} else if attr.Key == "name" {
					input.Name = attr.Val
				} else if attr.Key == "value" {
					input.Value = attr.Val
				}
			}
			mutexes.MutexFormInputs.Lock()
			content.FormInputs = append(content.FormInputs, input)
			mutexes.MutexFormInputs.Unlock()

		case "meta":
			var meta _MetaTag
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					meta.Name = attr.Val
				} else if attr.Key == "content" {
					meta.Content = attr.Val
				}
			}
			mutexes.MutexMetaTags.Lock()
			content.MetaTags = append(content.MetaTags, meta)
			mutexes.MutexMetaTags.Unlock()

		case "h1", "h2", "h3", "h4", "h5", "h6":
			var header _Header
			header.Tag = n.Data
			if n.FirstChild != nil {
				header.Text = n.FirstChild.Data
			}
			mutexes.MutexHeaders.Lock()
			content.Headers = append(content.Headers, header)
			mutexes.MutexHeaders.Unlock()
		}
	} else if n.Type == html.TextNode {
		if text := strings.TrimSpace(n.Data); text != "" {
			mutexes.MutexBodyInnerText.Lock();
			content.BodyInnerText += text + " "
			mutexes.MutexBodyInnerText.Unlock();
		}
	}
	// Recur for child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		wg.Add(1)
		go processNode(c, content, mutexes, wg)
	}
}

func parsePage(url string, body string) (PageData, error) {
	var wg sync.WaitGroup
	
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return PageData{}, err
	}
	wg.Add(1)
	var htmlContent ScrapedDefault
	var mutexes ScrapedDefaultMutex
	processNode(doc, &htmlContent, &mutexes, &wg)
	wg.Wait()
	return PageData {
		URL: url,
		HTMLContent: htmlContent,
	}, nil
}
