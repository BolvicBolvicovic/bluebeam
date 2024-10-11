package analyzer

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"net/http"
	"github.com/BolvicBolvicovic/bluebeam/criterias"
	"os/exec"
	"encoding/json"
)

type _Button struct {
	Text    string `json:"text"`
	OnClick string `json:"onclick"`
	ID      string `json:"id"`
	Classes string `json:"classes"`
}

type _Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type _Image struct {
	Src     string `json:"src"`
	Alt     string `json:"alt"`
	Classes string `json:"classes"`
}

type _FormInput struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type _MetaTag struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type _Header struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

type ScrapedDefault struct {
	Username     string       `json:"username"`
	SessionKey   string       `json:"sessionkey"`
	Links        []_Link      `json:"links"`
	Buttons      []_Button    `json:"buttons"`
	Images       []_Image     `json:"images"`
	FormInputs   []_FormInput `json:"formInputs"`
	MetaTags     []_MetaTag   `json:"metaTags"`
	Headers      []_Header    `json:"headers"`
	BodyInnerText string      `json:"bodyText"`
}

type LLMQuestions struct {
	SystemMessage	string `json:"systemmessage"`
	Data		ScrapedDefault `json:"data"`
	Features	[]criterias.Feature `json:"features"`
}

type LLMQuestion struct {
	SystemMessage	string `json:"systemmessage"`
	Data		ScrapedDefault `json:"data"`
	Feature 	criterias.Feature `json:"feature"`
}

type LLMResponse struct {
	response	string
	mutex		sync.Mutex
}

var wg sync.WaitGroup

func sendLLMQuestion(f criterias.Feature, sd *ScrapedDefault, r *LLMResponse) {
	defer wg.Done()

	question := LLMQuestion {
		SystemMessage: "You extract feature from data into JSON data if you find the feature in data else precise otherwise in the JSON data",
		Data: *sd,
		Feature: f,
	}
	questionJSON, err := json.Marshal(question)
	if err != nil {
		log.Println(err)
		return
	}
	var strResponse string
	response, err := exec.Command(
			"/venv/bin/python3",
			"analyzer/llm_client.py",
			string(questionJSON),
		).Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			strResponse = string(exitError.Stderr)
		} else {
			strResponse = err.Error()
		}
	} else {
		strResponse = string(response)
	}

	r.mutex.Lock()
	r.response += strResponse
	r.mutex.Unlock()
}

func Analyzer(c *gin.Context, sd ScrapedDefault) {
	crits, err := criterias.Get(c, sd.Username)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ("No criterias chosen or " + err.Error())})
		return
	}
	var response LLMResponse
	for _, feat := range crits.Features {
		wg.Add(1)
		go sendLLMQuestion(feat, &sd, &response)				
	}
	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"message": response.response})
}
