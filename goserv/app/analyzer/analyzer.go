package analyzer

import (
	"strings"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"net/http"
	"github.com/BolvicBolvicovic/bluebeam/criterias"
	"os/exec"
	"os"
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
	Links        []_Link      `json:"links"`
	Buttons      []_Button    `json:"buttons"`
	Images       []_Image     `json:"images"`
	FormInputs   []_FormInput `json:"formInputs"`
	MetaTags     []_MetaTag   `json:"metaTags"`
	Headers      []_Header    `json:"headers"`
	BodyInnerText string      `json:"bodyText"`
}

type ScrapedUrls struct {
	Urls        []string      `json:"urls"`
}

type LLMQuestions struct {
	SystemMessage	string `json:"systemmessage"`
	Data		json.RawMessage `json:"data"`
	Features	[]criterias.Feature `json:"features"`
}

type LLMQuestion struct {
	SystemMessage	string `json:"systemmessage"`
	Data		json.RawMessage `json:"data"`
	Feature 	criterias.Feature `json:"feature"`
}

type LLMResponse struct {
	Response	[]json.RawMessage
	mutex		sync.Mutex
}

type WebsiteGroup struct {
	Websites	[]Website
	mutex		sync.Mutex
}

func sendLLMQuestion(f criterias.Feature, sd *json.RawMessage, r *LLMResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	tempFile, err := os.CreateTemp("", "sd_data_*.json")
	if err != nil {
	    log.Println("Error creating temp file:", err)
	    return
	}
	defer os.Remove(tempFile.Name()) // Clean up after usage
	
	question := LLMQuestion {
		SystemMessage: "You extract feature from data into JSON data if you find the feature in data else precise otherwise in the JSON data",
		Data: *sd,
		Feature: f,
	}
	questionJSON, err := json.Marshal(question)
	if err != nil {
		log.Println("here",err)
		return
	}
	if _, err := tempFile.Write(questionJSON); err != nil {
	    log.Println("Error writing to temp file:", err)
	    return
	}


	var strResponse string
	response, err := exec.Command(
			"/venv/bin/python3",
			"analyzer/llm_client.py",
			tempFile.Name(),
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
	log.Println("LLMResponse:", strResponse)
	if (strings.Contains(strResponse, "error: ")) {
		return
	}
	jsonResponse := json.RawMessage(strResponse)

	r.mutex.Lock()
	r.Response = append(r.Response, jsonResponse)
	r.mutex.Unlock()
}

func Analyzer(c *gin.Context, sd ScrapedDefault, username string) {
	var wg sync.WaitGroup
	crits, index_file, err := criterias.Get(c, username)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ("No criterias chosen or " + err.Error())})
		return
	}
	if index_file == -1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No criterias chosen"})
		return
	}
	var response LLMResponse
	sdm, err := json.Marshal(sd)
	if err != nil {
		log.Println("here",err)
		return
	}
	rsdm := json.RawMessage(string(sdm))
	for _, feat := range crits[index_file].Features {
		wg.Add(1)
		go sendLLMQuestion(feat, &rsdm, &response, &wg)				
	}
	wg.Wait()
	var finalResponse [][]json.RawMessage
	finalResponse = append(finalResponse, response.Response)
	c.JSON(http.StatusOK, gin.H{"message": finalResponse})
}

func HandleUrls(c *gin.Context, su ScrapedUrls, username string) {
	var wgUrls sync.WaitGroup

	crits, index_file, err := criterias.Get(c, username)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ("No criterias chosen or " + err.Error())})
		return
	}
	if index_file == -1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No criterias chosen"})
		return
	}
	var crawledWebsites WebsiteGroup
	for _, item := range su.Urls {
		wgUrls.Add(1)
		go crawlWebsite(item, &crawledWebsites, &wgUrls)
	}
	wgUrls.Wait()
	
	responses := make([]LLMResponse,150)
	for i, site := range crawledWebsites.Websites {
		if i >= 150 {
			break
		}
		if site.Pages == nil {
			continue
		}
		wgUrls.Add(1)
		go func() {
			defer wgUrls.Done()
			var wg sync.WaitGroup
			sdm, err := json.Marshal(site.Pages)
			if err != nil {
				log.Println("here",err)
				return
			}
			rsdm := json.RawMessage(string(sdm))
			for _, feat := range crits[index_file].Features {
				wg.Add(1)
				go sendLLMQuestion(feat, &rsdm, &responses[i], &wg)				
			}
			wg.Wait()
		}()
	}
	wgUrls.Wait()
	var finalResponses [][]json.RawMessage
	for _, response := range responses {
		if response.Response == nil {
			break
		}
		finalResponses = append(finalResponses, response.Response)
	}
	c.JSON(http.StatusOK, gin.H{"message": finalResponses})
}
