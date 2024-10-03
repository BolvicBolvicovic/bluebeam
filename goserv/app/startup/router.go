package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
//	"time"
	"github.com/BolvicBolvicovic/scraper/api"
)

func BuildRouter() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = append(config.AllowOrigins, "https:/localhost/*")
	router.Use(cors.New(config))

	router.GET("/ping", api.Pong)
	router.POST("/login", api.Login)
	router.POST("/register_account", api.ResgisterAccount)
	router.POST("/analyse", api.Analyse)
	return router
}
