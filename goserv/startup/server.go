package startup

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/BolvicBolvicovic/scraper/config"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
)

type Shutdown = func()

var db *sql.DB

func Server() {
	env := config.NewEnv(".env", true)
	router, addr, shutdown := create(env)
	defer shutdown()
	fmt.Printf("Launching server on: %v:%v\n", env.ServerHost, env.ServerPort)
	http.ListenAndServe(addr, router)
}

func connectDB(env *config.Env) *sql.DB {
	db_addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		env.DBUser, env.DBUserPwd, env.DBHost, env.DBPort, env.DBName,
	)
	db, err := sql.Open("mysql", db_addr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	fmt.Println("Connecting to DB")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to DB")
	return db
}

func create(env *config.Env) (*gin.Engine, string, Shutdown) {
	
	db = connectDB(env)

	router := BuildRouter()
	
	addr := fmt.Sprintf("%s:%d", env.ServerHost, env.ServerPort)
	shutdown := func() {
		db.Close()
	}
	return router, addr, shutdown
}
