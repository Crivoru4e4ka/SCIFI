package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var err error

	connStr := "host=localhost port=5432 user=postgres password=btstop132013 dbname=SRA sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Cannot connect: ", err)
	}

	log.Println("Connected to PostgreSQL")
}
