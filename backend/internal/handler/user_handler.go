package handler

import (
	"backend/internal/model"
	"backend/internal/service"
	"encoding/json"
	"net/http"
)

// UserHandler handles HTTP requests related to user operations
type UserHandler struct {
	Service *service.UserService // Service layer dependency for business logic
}

// Register handles user registration requests
// It expects a JSON payload with user details in the request body
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Create an empty User struct to hold the decoded data
	var user model.User

	// Decode the JSON request body into the User struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// Return 400 Bad Request if JSON is malformed
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call the service layer to register the user
	if err := h.Service.RegisterUser(&user); err != nil {
		// Return 500 Internal Server Error if registration fails
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	// Set HTTP status to 201 Created for successful registration
	w.WriteHeader(http.StatusCreated)

	// Return success response as JSON
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registration successful",
	})
}