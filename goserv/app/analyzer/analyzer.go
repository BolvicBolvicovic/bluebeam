package analyzer

import (
	"strings"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"net/http"
	"github.com/BolvicBolvicovic/bluebeam/criterias"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"fmt"
	"errors"
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
	Urls		[]string  `json:"urls"`
	Ai		string	  `json:"ai"`
	Sanitizer	string	  `json:"sanitizer"`
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
	Websites	[][]PageData  `json:"crawledwebsites"`
	mutex		sync.Mutex
}

func getAPIKey(name string, username string) (string, error) {
	key := func() string {
		if name == "GEMINI_API_KEY" { return "gemini_api_key"}
		return "openai_api_key"
	}() 
	query := `
SELECT
	%s
FROM
	users
WHERE
	username = ?;
	`
	query = fmt.Sprintf(query, key)
	row := database.Db.QueryRow(query, username)

	var apiKey sql.NullString
	if err := row.Scan(&apiKey); err != nil {
		return "", err
	}
	if !apiKey.Valid {
		return "", errors.New("no key")
	}
	return fmt.Sprintf("%s=%s", name, apiKey.String), nil	
}

func sendLLMQuestion(f criterias.Feature, sd *json.RawMessage, r *LLMResponse, wg *sync.WaitGroup, ai string, ai_key string) {
	defer wg.Done()

	ai_client := ""
	switch ai {
	case "gpt-4o-mini":
		ai_client = "analyzer/openai/llm_client.py"
	case "gemini-1.5-flash", "gemini-1.5-pro":
		ai_client = "analyzer/gemini/llm_client.py"
	default:
	    log.Printf("AI %s does not exist\n", ai)
	    return
	}
	tempFile, err := os.CreateTemp("", "sd_data_*.json")
	if err != nil {
	    log.Println("Error creating temp file:", err)
	    return
	}
	defer os.Remove(tempFile.Name())
	
	question := LLMQuestion {
		SystemMessage: "You extract feature from data into the JSON output if feature is found in data else precise otherwise in the JSON output",
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
	command := exec.Command("/venv/bin/python3", ai_client, tempFile.Name(), ai)
	command.Env = append(command.Env, ai_key)
	response, err := command.Output()
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

func sanitizeCrawledWebsites(crawledWebsites WebsiteGroup, ai_name string, ai_key string, sanitizer string) (WebsiteGroup, error) {
	var finalResponse WebsiteGroup
	ai_client := ""
	switch ai_name {
	case "gpt-4o-mini":
		ai_client = "analyzer/openai/llm_sanitizer.py"
	case "gemini-1.5-flash", "gemini-1.5-pro":
		ai_client = "analyzer/gemini/llm_sanitizer.py"
	default:
		return crawledWebsites, errors.New(fmt.Sprintf("AI %s does not exist\n", ai_name))
	}
	tempFile, err := os.CreateTemp("", "cw_data_*.json")
	if err != nil {
		return crawledWebsites, err
	}
	defer os.Remove(tempFile.Name())
	type Question struct {
		CrawledWebsites [][]PageData	`json:"crawledwebsites"`
		AIName		string		`json:"ainame"`
		SystemMessage	string		`json:"systemmessage"`
		Sanitizer	string		`json:"sanitizer"`
	}
	question := Question{
		CrawledWebsites: crawledWebsites.Websites,
		AIName: ai_name,
		SystemMessage: `
You are a Data Sanitizer Assistant designed to process datasets, preserving their original structure while removing irrelevant or extraneous information. When a dataset is provided:

    Preserve Structure:
        Return the dataset in the exact structure and order it was received, including all keys, subkeys, and hierarchical relationships.
        Do not add, remove, or rename keys unless instructed explicitly in the user prompt.

    Sanitize Data Content:
        Identify and remove any irrelevant or extraneous information within the fields based on user-provided criteria.
        Ensure that each field contains only data directly relevant to the task, flagging irrelevant content for removal.

    Consistency and Accuracy:
        Retain valid data and correct minor errors within fields, such as format inconsistencies or minor typographical errors.
        Do not alter values or units unless explicitly requested.

    Return Format:
        Return the sanitized dataset without any additional comments or explanation, ensuring it matches the format of the original data exactly.
        If the dataset contains multiple levels (e.g., nested objects), apply sanitization consistently at each level.
	Do not return the sanitizater specification as it is the user instruction on how to sanitize.

Your task is to ensure the dataset returned is cleaned of irrelevant data and remains in a consistent, structured format that closely mirrors the input.
		`,
		Sanitizer: sanitizer,
	}

	marshaledQuestion, err := json.Marshal(question)
	if err != nil {
		return crawledWebsites, err
	}
	if _, err := tempFile.Write(marshaledQuestion); err != nil {
		return crawledWebsites, err
	}
	var strResponse string
	command := exec.Command("/venv/bin/python3", ai_client, tempFile.Name())
	command.Env = append(command.Env, ai_key)
	response, err := command.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			strResponse = string(exitError.Stderr)
		} else {
			strResponse = err.Error()
		}
	} else {
		strResponse = string(response)
	}
	if (strings.Contains(strResponse, "error: ")) {
		return crawledWebsites, errors.New(strResponse)
	}
	if err = json.Unmarshal([]byte(strResponse), finalResponse.Websites); err != nil {
		log.Printf("The error is here:", strResponse)
		return crawledWebsites, err
	}
	
	return finalResponse, nil
}

func Analyzer(c *gin.Context, sd ScrapedDefault, username string) {
	var wg sync.WaitGroup
	ai_key, err := getAPIKey("OPENAI_API_KEY", username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No API key for openai"})
		return
	}
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
		go sendLLMQuestion(feat, &rsdm, &response, &wg, "gpt-4o-mini", ai_key)
	}
	wg.Wait()
	var finalResponse [][]json.RawMessage
	finalResponse = append(finalResponse, response.Response)
	c.JSON(http.StatusOK, gin.H{"message": finalResponse})
}

func HandleUrls(c *gin.Context, su ScrapedUrls, username string) {
	var wgUrls sync.WaitGroup
	ai_name := func () string {
		if strings.Contains(su.Ai, "gemini") { return "GEMINI_API_KEY" }
		return "OPENAI_API_KEY"
	}()
	ai_key, err := getAPIKey(ai_name, username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ("No API key for the chosen AI" + err.Error())})
		return
	}

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
	
	if su.Sanitizer != "" {
		crawledWebsites, err = sanitizeCrawledWebsites(crawledWebsites, su.Ai, ai_key, su.Sanitizer)
		if err != nil {
			log.Println("error sanitizer:", err)
		} else {
			log.Println("first sanitized page:", crawledWebsites.Websites[0][0])
		}
	}	
	
	responses := make([]LLMResponse,150)
	for i, site := range crawledWebsites.Websites {
		if i >= 150 {
			break
		}
		if site == nil {
			continue
		}
		wgUrls.Add(1)
		go func() {
			defer wgUrls.Done()
			var wg sync.WaitGroup
			sdm, err := json.Marshal(site)
			if err != nil {
				log.Println("here",err)
				return
			}
			rsdm := json.RawMessage(string(sdm))
			for _, feat := range crits[index_file].Features {
				wg.Add(1)
				go sendLLMQuestion(feat, &rsdm, &responses[i], &wg, su.Ai, ai_key)
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
