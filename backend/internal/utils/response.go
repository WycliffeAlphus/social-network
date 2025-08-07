package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse structure for consistent API responses
type JSONResponse struct {
	Message string      `json:"message"`
	Data    any `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RespondWithJSON sends JSON respones
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response := JSONResponse{
		Message: http.StatusText(status),
		Data:    payload,
	}
	if status >= 400 {
		response.Error = http.StatusText(status)
		response.Data = nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, status int, message string) {
	response := JSONResponse{
		Message: message,
		Error:   message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}



