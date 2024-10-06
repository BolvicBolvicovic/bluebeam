package criterias

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"sync"
)

var (
	aeadInstance tink.AEAD
	once         sync.Once
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		log.Println(err)
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type Feature struct {
	Topic 		string    `json:"topic" ` 
	FeatureName 	string    `json:"featurename" `
	FeatureType	string	  `json:"featuretype" `
	Description 	string    `json:"description" `
	MinimumCriterias []string `json:"minimumcriterias" `
	YesCases 	[]string  `json:"yescases" `
	NoCases 	[]string  `json:"nocases" `
}

type Criterias struct {
	Username	string	  `json:"username"`
	SessionKey	string    `json:"sessionkey"`
	Features	[]Feature `json:"features"`
}

func Store(c *gin.Context, crits Criterias) {
	log.Println(crits)
	data, _ := json.Marshal(crits)
	encryptedData, err := aeadInstance.Encrypt(data, nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	//TODO: Resquest to db to store data
	c.JSON(http.StatusOK, gin.H{"message": "Criterias well recieved, Data processed!"})
}

func Get(c *gin.Context, username string) Criterias, error {
	//TODO: Resquest to db to get data
	decryptedData, err := aeadInstance.Decrypt(encryptedData, nil)
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}
	var decryptedStruct Criterias
	err = json.Unmarshal(decryptedData, &decryptedStruct)
	if err != nil {
		log.Fatalf("Failed to deserialize decrypted data: %v", err)
	}
	log.Println("Criterias decrypted:", decryptedStruct)
}

func SetKey() {
	once.Do(func() {
		kh, err := keyset.NewHandle(aead.AES256CTRHMACSHA256KeyTemplate())
		if err != nil {
			log.Fatalf("Failed to generate new key handle: %v", err)
		}

		a, err := aead.New(kh)
		if err != nil {
			log.Fatalf("Failed to get AEAD primitive: %v", err)
		}

		aeadInstance = a
	})
}
