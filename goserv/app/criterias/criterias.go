package criterias

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/firdasafridi/gocrypt"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

var aeskey string

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		log.Println(err)
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type Feature struct {
	Topic 		string    `json:"topic" gocrypt:"aes"` 
	FeatureName 	string    `json:"featurename" gocrypt:"aes"`
	FeatureType	string	  `json:"featuretype" gocrypt:"aes"`
	Description 	string    `json:"description" gocrypt:"aes"`
	MinimumCriterias []string `json:"minimumcriterias" gocrypt:"aes"`
	YesCases 	[]string  `json:"yescases" gocrypt:"aes"`
	NoCases 	[]string  `json:"nocases" gocrypt:"aes"`
}

type Criterias struct {
	Username	string	  `json:"username"`
	SessionKey	string    `json:"sessionkey"`
	Features	[]Feature `json:"features"`
}

func Store(c *gin.Context, crits Criterias) {
	criterias := Criterias{}

	if err := c.ShouldBindJSON(&criterias); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	aesOpt, err := gocrypt.NewAESOpt(aeskey)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	log.Println(criterias)

	cryptRunner := gocrypt.New(&gocrypt.Option {AESOpt : aesOpt,})
	if err = cryptRunner.Encrypt(&criterias); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	log.Println("Criterias encrypted:", criterias)
}

func SetKey() {
	aeskey, _ = randomHex(20)
	aeskey = strings.ToLower(aeskey)
}
