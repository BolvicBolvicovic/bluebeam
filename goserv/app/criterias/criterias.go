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
	"github.com/google/tink/go/proto/tink_go_proto"
	"github.com/google/tink/go/proto/aes_gcm_go_proto"
	"github.com/golang/protobuf/proto"
	"encoding/json"
)

var aeadInstance tink.AEAD

type Feature struct {
	Topic 		string  `json:"topic" ` 
	FeatureName 	string  `json:"featurename" `
	FeatureType	string	`json:"featuretype" `
	Description 	string  `json:"description" `
	MinimumCriterias string `json:"minimumcriterias" `
	YesCases 	string  `json:"yescases" `
	NoCases 	string  `json:"nocases" `
	AuthorizedData	string  `json:"authorizeddata"`
}

type Criterias struct {
	Features	[]Feature `json:"features"`
	FileName	string	  `json:"filename"`	
}

func Store(c *gin.Context, crits Criterias, username string) {
	currentCriteriasFiles,_, _ := Get(c, username)
	currentCriteriasFiles = append(currentCriteriasFiles, crits)
	data, err := json.Marshal(currentCriteriasFiles)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	encryptedData, err := aeadInstance.Encrypt(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	query := `
UPDATE
	users
SET
	criterias_files = ?,
	current_file_index = ?
WHERE
	username = ?;
	`
	if _, err := database.Db.Exec(query, encryptedData, len(currentCriteriasFiles) - 1, username); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading file into database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Criterias well recieved, Data processed!"})
}

func Get(c *gin.Context, username string) ([]Criterias, int, error) {
	query := `
SELECT
	criterias_files,
	current_file_index
FROM
	users
WHERE
	username = ?;
	`
	row := database.Db.QueryRow(query, username)
	var encryptedData sql.Null[[]byte]
	var index_file int
	if err := row.Scan(&encryptedData, &index_file); err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return make([]Criterias, 0), -1, err
	}
	if !encryptedData.Valid {
		return make([]Criterias, 0), -1, nil
	}
	decryptedData, err := aeadInstance.Decrypt(encryptedData.V, nil)
	if err != nil {
		return make([]Criterias, 0), -1, err
	}
	decryptedStruct := make([]Criterias, 0)
	err = json.Unmarshal(decryptedData, &decryptedStruct)
	if err != nil {
		return make([]Criterias, 0), -1, err
	}
	return decryptedStruct, index_file, nil
}

func SetKey() {
	var kek = aes_gcm_go_proto.AesGcmKeyFormat {
		Version: 0,
		KeySize: 32,
	}
	val, err := proto.Marshal(&kek)
    	if err != nil {
		log.Fatalf("Failed to generate new kek handle: %v", err)
    	}
	keyTemplate := tink_go_proto.KeyTemplate{
		TypeUrl: "type.googleapis.com/google.crypto.tink.AesGcmKey",
		Value:   val,
		OutputPrefixType: tink_go_proto.OutputPrefixType_TINK,
	}

    	kekh, err := keyset.NewHandle(&keyTemplate)
    	if err != nil {
		log.Fatalf("Failed to generate new kek handle: %v", err)
    	}

    	kekInstance, err := aead.New(kekh)
    	if err != nil {
		log.Fatalf("Failed to generate new kek handle: %v", err)
    	}

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
		if err := kh.Write(writer, kekInstance); err != nil {
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
	kh, err := keyset.Read(reader, kekInstance)
	if err != nil {
		log.Fatalf("Failed to deserialize key handle: %v", err)
	}

	aeadInstance, err = aead.New(kh)
	if err != nil {
		log.Fatalf("Failed to create AEAD instance: %v", err)
	}
}
