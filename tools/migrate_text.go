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
	sql, _ := os.ReadFile("migrations/2026-05-27_dataset_metadata_text.sql")
	_, err := db.DB.Exec(string(sql))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migration applied")
}
