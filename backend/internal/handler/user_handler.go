package handler

import (
	"backend/internal/model"
	"backend/internal/service"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Service.RegisterUser(&user); err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}
