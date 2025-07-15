package main

import (
	"backend/internal/routes"
	"backend/pkg/db/sqlite"
	"log"
	"net/http"
)

func main() {
	// Connect to the SQLite database and apply migrations
	db, err := sqlite.ConnectAndMigrate()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Register all routes (handlers)
	routes.RegisterRoutes(db)

	// Serve uploaded files from the /uploads/ directory
	// This allows accessing files at http://localhost:8080/uploads/<filename>
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server started on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
