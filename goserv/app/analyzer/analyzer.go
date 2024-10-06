package analyzer

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"net/http"
	"github.com/BolvicBolvicovic/bluebeam/criterias"
)

type _Buttons struct {
	Text	string `json:"text"`
	OnClick string `json:"onclick"`
}

type ScrapedDefault struct {
	Username	string `json:username`
	SessionKey	string `json:"sessionkey"`
	Links		[]string `json:"links"`
	Buttons		[]_Buttons `json:"buttons"`
	PageHtml	string `json:"pageHtml"`
}

var wg sync.WaitGroup

func checkLinks(links []string) {
	defer wg.Done()
	for i, link := range links {
		log.Println("link", i + 1, link)
	}
}

func checkButtons(buttons []_Buttons) {
	defer wg.Done()
	for i, button := range buttons {
		log.Println("button", i + 1, button)
	}
}

func checkHTML(html string) {
	defer wg.Done()
	log.Println("html", html)
}

func Analyzer(c *gin.Context, sd ScrapedDefault) {
	crits, err := criterias.Get(c, sd.Username)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	wg.Add(3)
	go checkLinks(sd.Links)
	go checkButtons(sd.Buttons)
	go checkHTML(sd.PageHtml)
	wg.Wait()
	log.Println(crits)
	c.JSON(http.StatusOK, gin.H{"message": "Page well recieved, Data processed!"})
}
