package analyser

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"sync"
//	"github.com/BolvicBolvicovic/scraper/database"
//	"database/sql"
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

var upgrader = websocket.Upgrader {
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func Analyser(c *gin.Context, sd ScrapedDefault) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()
	wg.Add(3)
	go checkLinks(sd.Links)
	go checkButtons(sd.Buttons)
	go checkHTML(sd.PageHtml)
	wg.Wait()
	ws.WriteMessage(websocket.TextMessage, []byte("Analyse is a success!"))
}
