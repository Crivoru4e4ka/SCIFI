package main

import (
	"log"
	"net/http"

	"project-MVP/db"
	"project-MVP/routes"
	"project-MVP/services"

	"github.com/joho/godotenv"
)

// @title Scifi API
// @version 1.0
// @description API сервера для системы управления научными проектами Scifi.
// @host localhost:8080
// @BasePath /
func main() {
	// Загружаем переменные окружения из .env файла
	// Если файл не найден или это продакшн сервер – это нормально
	errErr := godotenv.Load()
	if errErr != nil {
		log.Println("Warning: Could not load .env file (OK for production with system env vars)")
	}

	db.Connect()

	services.InitDefaultStores()

	services.InitRoles()
	services.InitPermissions()
	services.InitRolePermissions()

	r := routes.InitRoutes()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
