package database

import (
	"github.com/BolvicBolvicovic/scraper/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
)
type Shutdown func()
var Db *sql.DB

func ConnectDB(env *config.Env) Shutdown {
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
	if Db == nil {
		log.Fatal("Error, DB not Connected")
	}
	return func() {
		Db.Close()
	}
}
