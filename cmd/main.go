package main

import (
	"log"

	"github.com/BuddhiLW/DL-backend/internal/db"
	"github.com/BuddhiLW/DL-backend/internal/handlers"
)

func main() {
	// Initialize the database
	database := db.InitDatabase()

	// Set up the router
	router := handlers.SetupRouter(database)

	// Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}
