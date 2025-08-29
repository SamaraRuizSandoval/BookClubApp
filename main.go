package main

import (
	"log"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/app"
)

func main() {
	// Create default server configuration
	config := app.DefaultConfig()

	// Create new server instance
	server := app.NewServer(config)

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}
