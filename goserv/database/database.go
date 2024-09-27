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
	Db = db
	if Db == nil {
		log.Fatal("Error, DB not Connected")
	}
	if _, err := Db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(100) NOT NULL UNIQUE,
	password VARCHAR(100) NOT NULL,
	session_key VARCHAR(32),
	creation_key_time VARCHAR(50)
);
	`); err != nil {
		log.Fatal(err)
	}

	return func() {
		Db.Close()
	}
}
