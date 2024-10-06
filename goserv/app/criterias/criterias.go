package criterias

import (
	"encoding/base64"
	"bytes"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
	"encoding/json"
)

var aeadInstance tink.AEAD

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
	query := `
UPDATE
	users
SET
	criterias_file = ?
WHERE
	username = ?;
	`
	if _, err := database.Db.Exec(query, encryptedData, crits.Username); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading file into database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Criterias well recieved, Data processed!"})
}

func Get(c *gin.Context, username string) (Criterias, error) {
	query := `
SELECT
	criterias_file
FROM
	users
WHERE
	username = ?;
	`
	row := database.Db.QueryRow(query, username)
	var encryptedData sql.Null[[]byte]
	if err := row.Scan(&encryptedData); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return Criterias{}, err
	}
	decryptedData, err := aeadInstance.Decrypt(encryptedData.V, nil)
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}
	var decryptedStruct Criterias
	err = json.Unmarshal(decryptedData, &decryptedStruct)
	if err != nil {
		log.Fatalf("Failed to deserialize decrypted data: %v", err)
	}
	return decryptedStruct, nil
}

func SetKey() {
	query := `
SELECT
	first_key
FROM
	decrypt_keys
WHERE
	id = 0;
	`
	row := database.Db.QueryRow(query)

	var keyBase64 sql.NullString
	if err := row.Scan(&keyBase64); err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}

		kh, err := keyset.NewHandle(aead.AES256CTRHMACSHA256KeyTemplate())
		if err != nil {
			log.Fatalf("Failed to generate new key handle: %v", err)
		}

		buf := new(bytes.Buffer)
		writer := keyset.NewBinaryWriter(buf)
		if err := kh.Write(writer, nil); err != nil {
			log.Fatalf("Failed to serialize key handle: %v", err)
		}
		keyData := buf.Bytes()

		keyBase64Str := base64.StdEncoding.EncodeToString(keyData)

		insertQuery := `
INSERT INTO
	decrypt_keys
	(first_key)
VALUES
	(?);
		`
		if _, err := database.Db.Exec(insertQuery, keyBase64Str); err != nil {
			log.Fatal(err)
		}

		aeadInstance, err = aead.New(kh)
		if err != nil {
			log.Fatalf("Failed to create AEAD instance: %v", err)
		}
		return
	}

	keyData, err := base64.StdEncoding.DecodeString(keyBase64.String)
	if err != nil {
		log.Fatalf("Failed to decode key material: %v", err)
	}

	reader := keyset.NewBinaryReader(bytes.NewReader(keyData))
	kh, err := keyset.Read(reader, nil)
	if err != nil {
		log.Fatalf("Failed to deserialize key handle: %v", err)
	}

	aeadInstance, err = aead.New(kh)
	if err != nil {
		log.Fatalf("Failed to create AEAD instance: %v", err)
	}
}
