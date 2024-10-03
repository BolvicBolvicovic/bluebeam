package api

import (
	"github.com/gin-gonic/gin"
)

func Pong(ctx *gin.Context) {
	ctx.JSON(200, gin.H { "message": "pong", })
}
