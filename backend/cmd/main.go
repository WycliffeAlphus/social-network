package main

import (
	"backend/internal/routes"
	"backend/pkg/db/sqlite"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := sqlite.ConnectAndMigrate()
	if err != nil {
		fmt.Println(err.Error())
	}

	routes.RegisterRoutes(db)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
