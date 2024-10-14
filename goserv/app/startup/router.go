package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/BolvicBolvicovic/bluebeam/api"
	"log"
	"time"
	"io/ioutil"
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

func SheetsServiceMiddleware() gin.HandlerFunc {
	ctx := context.Background()
	creds, err := ioutil.ReadFile("startup/googlecredentials.json")
	if err != nil {
		log.Fatal("Unable to read credentials file:", err)
	}
	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope, drive.DriveFileScope)
	if err != nil {
		log.Fatal("Unable to parse credentials file:", err)
	}
	client := config.Client(ctx)
	sService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal("Unable to retrieve Sheets client:", err)
	}
	dService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal("Unable to retrieve Drive client:", err)
	}
	return func(c *gin.Context) {
		c.Set("sheetsService", sService)
		c.Set("driveService", dService)
		c.Next()
	}
}

func BuildRouter() *gin.Engine {
	router := gin.Default()

	config := cors.Config {
		AllowOrigins: []string{"https://localhost", "moz-extension://"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
		ExposeHeaders:[]string{"Content-Type"},
		AllowWildcard:true,
		AllowBrowserExtensions: true,
		MaxAge: time.Hour * 12,
	}
	
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	router.Use(cors.New(config))
	router.Use(SheetsServiceMiddleware())
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", api.Pong)
	router.GET("/settings", api.Settings)
	router.GET("/initOAuth", api.InitOAuth)
	router.GET("/selectGoogleFile", api.SelectGoogleFile)

	router.POST("/login", api.Login)
	router.POST("/registerAccount", api.ResgisterAccount)
	router.POST("/analyze", api.Analyze)
	router.POST("/criterias", api.StoreCriterias)
	router.POST("/updateEmail", api.UpdateEmail)
	router.POST("/outputGoogleSpreadsheet", api.OutputGoogleSpreadsheet)
	router.POST("/processGoogleSpreadsheet", api.ProcessGoogleSpreadsheet)

	return router
}
