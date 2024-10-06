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
	aesOpt, err := gocrypt.NewAESOpt(aeskey)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	log.Println(crits)

	cryptRunner := gocrypt.New(&gocrypt.Option {AESOpt : aesOpt,})
	if err = cryptRunner.Encrypt(&crits); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	log.Println("Criterias encrypted:", crits)
	c.JSON(http.StatusOK, gin.H{"message": "Page well recieved, Data processed!"})
}

func SetKey() {
	aeskey, _ = randomHex(64)
	aeskey = strings.ToLower(aeskey)
}
