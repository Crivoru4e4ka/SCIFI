package main

import (
	"log"
	"net/http"

	"project-MVP/db"
	"project-MVP/routes"
	"project-MVP/services"

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

	services.InitRoles()
	services.InitPermissions()
	services.InitRolePermissions()

	r := routes.InitRoutes()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
