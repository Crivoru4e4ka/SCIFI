package main

import (
	"log"
	"net/http"

	"project-MVP/db"
	"project-MVP/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env файла
	// Если файл не найден или это продакшн сервер – это нормально
	errErr := godotenv.Load()
	if errErr != nil {
		log.Println("Warning: Could not load .env file (OK for production with system env vars)")
	}

	db.Connect()

	r := routes.InitRoutes()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
