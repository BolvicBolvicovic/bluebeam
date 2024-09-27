package startup

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/scraper/config"
	"github.com/BolvicBolvicovic/scraper/database"
	"fmt"
)

type Shutdown = func()

func Server() {
	env := config.NewEnv(".env", true)
	router, addr, shutdown := create(env)
	defer shutdown()
	fmt.Printf("Launching server on: %v:%v\n", env.ServerHost, env.ServerPort)
	http.ListenAndServe(addr, router)
}


func create(env *config.Env) (*gin.Engine, string, Shutdown) {
	
	shutdown := database.ConnectDB(env)

	router := BuildRouter()
	
	addr := fmt.Sprintf("%s:%d", env.ServerHost, env.ServerPort)
	return router, addr, shutdown
}
