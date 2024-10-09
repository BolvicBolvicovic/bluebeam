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

type _Buttons struct {
	Text	string `json:"text"`
	OnClick string `json:"onclick"`
}

type ScrapedDefault struct {
	Username	string `json:username`
	SessionKey	string `json:"sessionkey"`
	Links		[]string `json:"links"`
	Buttons		[]_Buttons `json:"buttons"`
	BodyInnerText	string `json:"bodyText"`
}

type LLMQuestions struct {
	SystemMessage	string `json:"systemmessage"`
	Data		ScrapedDefault `json:"data"`
	Features	[]criterias.Feature `json:"features"`
}

type LLMQuestion struct {
	systemMessage	string `json:"systemmessage"`
	data		ScrapedDefault `json:"data"`
	feature 	criterias.Feature `json:"feature"`
}

type LLMResponse struct {
	response	string
	mutex		sync.Mutex
}

var wg sync.WaitGroup

func sendLLMQuestion(f criterias.Feature, sd *ScrapedDefault, r *LLMResponse) {
	defer wg.Done()

	question := LLMQuestion {
		systemMessage: "You extract feature from data into JSON data if you find the feature in data else precise otherwise in the JSON data",
		data: *sd,
		feature: f,
	}
	questionJSON, err := json.Marshal(question)
	if err != nil {
		log.Println(err)
		return
	}

	response, err := exec.Command("python3", "analyzer/llm_client.py", string(questionJSON)).CombinedOutput()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(response)

	r.mutex.Lock()
	r.response += string(response)
	r.mutex.Unlock()
}

func Analyzer(c *gin.Context, sd ScrapedDefault) {
	crits, err := criterias.Get(c, sd.Username)
	if err != nil && err != sql.ErrNoRows {
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
