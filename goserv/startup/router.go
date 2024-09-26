package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/scraper/api"
)

func BuildRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/pong", api.Pong)
	return router
}
