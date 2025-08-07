package handler

import (
	"backend/internal/model"
	"backend/internal/service"
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
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
	imgURL := ""
	imageUrl, err := utils.HandlePostImageUpload(r, maxUploadSize, "avatarImage")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if imageUrl.Valid {
		imgURL = imageUrl.String
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
