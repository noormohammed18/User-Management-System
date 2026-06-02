package main

import (
	"fmt"
	"log"
	"net/http"
	"user-management/config"
	"user-management/routes"
)

func main() {
	config.ConnectDB()
	routes.SetupRoutes()

	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}
