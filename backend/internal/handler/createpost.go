package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/internal/service"
	"backend/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PostCreationErrors struct {
	Title            string `json:"titleerror,omitempty"`
	Content          string `json:"contenterror,omitempty"`
	PostPrivacy      string `json:"privacyerror,omitempty"`
	PostImage        string `json:"imageerror,omitempty"`
	AllowedFollowers string `json:"followerserror,omitempty"`
}

const (
	MinTitleLength   = 7
	MaxTitleLength   = 77
	MinContentLength = 21
	MaxContentLength = 777
	maxUploadSize    = 20 * 1024 * 1024 // 20MB
)

func CreatePost(db *sql.DB, notificationService *service.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID := context.MustGetUser(r.Context()).ID
		postErrors := &PostCreationErrors{}

		post := model.Post{
			Id:         uuid.New().String(),
			UserId:     currentUserID,
			Title:      r.FormValue("title"),
			Content:    r.FormValue("content"),
			Visibility: r.FormValue("postPrivacy"),
			CreatedAt:  time.Now(),
		}

		// Check for group_id and update the post model
		if groupID := r.FormValue("group_id"); groupID != "" {
			post.GroupId = sql.NullString{String: groupID, Valid: true}
		}

		if post.Visibility == "" {
			post.Visibility = "public"
		}

		if allowedJSON := r.FormValue("allowedFollowers"); allowedJSON != "" {
			err := json.Unmarshal([]byte(allowedJSON), &post.AllowedFollowers)
			if err != nil {
				http.Error(w, "Invalid allowedFollowers format", http.StatusBadRequest)
				return
			}
		}

		//  handle image upload
		imageUrl, err := utils.HandlePostImageUpload(r, maxUploadSize, "postImage")
		if err != nil {
			postErrors.PostImage = err.Error()
			http.Error(w, postErrors.PostImage, http.StatusBadRequest)
			return
		}
		if imageUrl.Valid {
			post.ImageUrl = imageUrl
		}

		postErrors, hasErrors := validatePost(post)
		if hasErrors {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(postErrors)
			return
		}

		_, postInsertErr := db.Exec(`
			INSERT INTO posts (id, user_id, title, content, visibility, post_image, created_at, group_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			post.Id, post.UserId, post.Title, post.Content, post.Visibility, post.ImageUrl, post.CreatedAt, post.GroupId)
		if postInsertErr != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		if post.Visibility == "private" && len(post.AllowedFollowers) > 0 {
			for _, followerID := range post.AllowedFollowers {
				_, err := db.Exec(`
					INSERT INTO private_posts (post_id, user_id, created_at)
					VALUES (?, ?, ?)`, post.Id, followerID, post.CreatedAt)
				if err != nil {
					http.Error(w, "Failed to add private post access", http.StatusInternalServerError)
					return
				}
			}
		}

		// Notify followers
		if err := notificationService.CreatePostNotification(post.UserId, post.Id, nil); err != nil {
			log.Printf("Error creating post notification: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	}
}

func validatePost(post model.Post) (*PostCreationErrors, bool) {
	errors := &PostCreationErrors{}

	if len(post.Title) > MaxTitleLength {
		errors.Title = fmt.Sprintf("Title length too long. Keep it at %d max", MaxTitleLength)
	}
	if len(post.Title) < MinTitleLength {
		errors.Title = fmt.Sprintf("Title length too short. Keep it at least %d", MinTitleLength)
	}
	if len(post.Content) > MaxContentLength {
		errors.Content = fmt.Sprintf("Content length too long. Keep it at %d max", MaxContentLength)
	}
	if len(post.Content) < MinContentLength {
		errors.Content = fmt.Sprintf("Content length too short. Keep it at least %d", MinContentLength)
	}

	visibility := strings.ToLower(post.Visibility)
	if visibility != "public" && visibility != "private" && visibility != "almostprivate" && visibility != "group" {
		errors.PostPrivacy = "Invalid privacy value. Must be 'public', 'almost private', or 'private'"
	}

	if visibility == "private" && len(post.AllowedFollowers) == 0 {
		errors.AllowedFollowers = "Please select at least one follower for private posts"
	}

	hasErrors := errors.HasErrors()
	return errors, hasErrors
}

func (pe *PostCreationErrors) HasErrors() bool {
	return pe.Title != "" ||
		pe.Content != "" ||
		pe.PostPrivacy != "" ||
		pe.PostImage != "" ||
		pe.AllowedFollowers != ""
}
