package api

import (
	"fmt"
	"io"
	"net/http"
	"github.com/gin-gonic/gin"
)

func Pong(ctx *gin.Context) {
	ctx.JSON(200, gin.H { "message": "pong", })
}


func HandleDB(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	fmt.Println(url.Hostname())		
	io.WriteString(w, "It's Me\n")
}
