package config

import (
	"database/sql"
	"log"
)

var database *sql.DB

func InitDatabase(dbUrl string) {
	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	database = db

	log.Println("Database connected on " + dbUrl)
}

func GetDatabase() *sql.DB {
	return database
}

func CloseDatabase() {
	database.Close()
}
