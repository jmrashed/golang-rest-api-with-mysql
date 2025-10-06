package main

import (
	"log"

	"jmrashed/apps/userApp/route"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	route.SetupRoutes()
}
