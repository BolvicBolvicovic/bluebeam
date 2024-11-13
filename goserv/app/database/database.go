package database

import (
	"github.com/BolvicBolvicovic/bluebeam/config"
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
	return createTables()
}

func createTables() Shutdown {
	if _, err := Db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(100) NOT NULL UNIQUE,
	password VARCHAR(100) NOT NULL,
	email VARCHAR(50) NOT NULL,
	session_key VARCHAR(72),
	creation_key_time VARCHAR(72),
	output_files_ids VARBINARY(15000),
	criterias_files BLOB(150000),
	current_file_index INT,
	gemini_api_key VARCHAR(100),
	openai_api_key VARCHAR(200)

);
	`); err != nil {
		log.Fatal(err)
	}
	if _, err := Db.Exec(`
CREATE TABLE IF NOT EXISTS decrypt_keys (
	id INT AUTO_INCREMENT PRIMARY KEY,
	first_key VARBINARY(400)
);
	`); err != nil {
		log.Fatal(err)
	}
	return func() {
		Db.Close()
	}
}
