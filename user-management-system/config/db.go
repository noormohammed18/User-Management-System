package config

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

func ConnectDB() {

    connString := "server=localhost;database=userdb;trusted_connection=yes"

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db

	log.Println("Database connected")
}