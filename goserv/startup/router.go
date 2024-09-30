package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/scraper/api"
)

func BuildRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", api.Pong)
	router.POST("/login", api.Login)
	router.POST("/register_account", api.ResgisterAccount)
	router.POST("/analyse", api.Analyse)
	return router
}
