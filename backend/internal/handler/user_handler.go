package handler

import (
	"backend/internal/model"
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP operations
type UserHandler struct {
	Service *service.UserService
}

// Register handles user registration via multipart/form-data.
// It expects an image file under the field "avatar" and other user details as form values.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register request received")

	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Parse multipart form (10MB max memory)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Invalid form data",
		})
		return
	}

	// Handle avatar upload
	imgURL := ""
	file, handler, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()

		mimeType := handler.Header.Get("Content-Type")
		if !strings.HasPrefix(mimeType, "image/") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "error",
				"message": "Only image uploads allowed",
			})
			return
		}

		filename := uuid.New().String() + filepath.Ext(handler.Filename)
		filePath := "./uploads/" + filename

		dst, err := os.Create(filePath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "error",
				"message": "Failed to save image",
			})
			return
		}
		defer dst.Close()
		io.Copy(dst, file)

		imgURL = "/uploads/" + filename
	}

	// Parse fields
	user := model.User{
		ID:                uuid.New().String(),
		Email:             r.FormValue("email"),
		FirstName:         r.FormValue("firstName"),
		LastName:          r.FormValue("lastName"),
		Nickname:          r.FormValue("nickname"),
		About:             r.FormValue("aboutMe"),
		Password:          r.FormValue("password"),
		ProfileVisibility: r.FormValue("profileVisibility"),
		ImgURL:            imgURL,
		CreatedAt:         time.Now(),
	}

	// Parse date of birth
	dobStr := r.FormValue("dateOfBirth")
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Invalid date format. Use YYYY-MM-DD.",
		})
		return
	}
	user.DOB = dob

	// Call service to save user
	validationErrors, err := h.Service.RegisterUser(&user)
	if validationErrors != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErrors)
		return
	}

	if err != nil {
		fmt.Println("DB Error:", err) // log it for devs
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Registration failed",
			"error":   err.Error(),
		})
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Registration successful",
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"firstName":  user.FirstName,
			"lastName":   user.LastName,
			"nickname":   user.Nickname,
			"imgUrl":     user.ImgURL,
			"visibility": user.ProfileVisibility,
		},
	})
}
