package routes

import (
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"
	"database/sql"
	"net/http"
)

func RegisterRoutes(db *sql.DB) {
	userRepo := &repository.UserRepository{DB: db}
	userService := &service.UserService{Repo: userRepo}
	userHandler := &handler.UserHandler{Service: userService}

	http.HandleFunc("/register", userHandler.Register)
}
