package main

import (
	"log"
	"net/http"
	"project-MVP/db"
	"project-MVP/routes"
)

func main() {
	db.Connect()

	r := routes.InitRoutes()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
