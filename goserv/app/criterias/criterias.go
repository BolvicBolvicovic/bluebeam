package criterias

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Criterias struct {
	Topic 		string    `json:"topic"` 
	FeatureName 	string    `json:"featurename"`
	Description 	string    `json:"description"`
	MinimumCriterias []string `json:"minimumcriterias"`
	YesCases 	[]string  `json:"yescases"`
	NoCases 	[]string  `json:"nocases"`
}

func Store(c *gin.Context) {
	criterias := Criterias{}

	if err := c.ShouldBindJSON(&criterias); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Println(criterias)
}
