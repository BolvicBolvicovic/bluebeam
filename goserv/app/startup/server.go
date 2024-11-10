package startup

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/bluebeam/config"
	"github.com/BolvicBolvicovic/bluebeam/database"
	"github.com/BolvicBolvicovic/bluebeam/criterias"
	"fmt"
	"os"
)

type Shutdown = func()

func Server() {
	env := config.NewEnv(".env", true)
	router, addr, shutdown := create(env)
	defer shutdown()
	fmt.Printf("Launching server on: %v:%v\n", env.ServerHost, env.ServerPort)
	criterias.SetKey()
	err := http.ListenAndServeTLS(addr, os.Getenv("SSL_CERT"), os.Getenv("SSL_KEY"), router)
	fmt.Printf("Server shutting down... error:", err)
}


func create(env *config.Env) (*gin.Engine, string, Shutdown) {
	
	shutdown := database.ConnectDB(env)

	router := BuildRouter()
	
	addr := fmt.Sprintf("%s:%d", env.ServerHost, env.ServerPort)
	return router, addr, shutdown
}
