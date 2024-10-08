package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/BolvicBolvicovic/bluebeam/api"
	"log"
	"time"
)

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
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", api.Pong)
	router.GET("/settings", api.Settings)

	router.POST("/login", api.Login)
	router.POST("/register_account", api.ResgisterAccount)
	router.POST("/analyze", api.Analyze)
	router.POST("/criterias", api.StoreCriterias)

	return router
}
