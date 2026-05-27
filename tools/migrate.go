package main

import (
	"log"
	"os"
	"project-MVP/db"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db.Connect()

	sql, err := os.ReadFile("migrations/2026-05-27_dataset_lineage.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.DB.Exec(string(sql))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migration applied successfully")
}
