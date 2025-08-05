package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type PostCreationErrors struct {
	Title       string `json:"titleerror,omitempty"`
	Content     string `json:"contenterror,omitempty"`
	PostPrivacy string `json:"privacyerror,omitempty"`
	PostImage   string `json:"imageerror,omitempty"`
	AllowedFollowers string `json:"followerserror,omitempty"`
}

const (
	MinTitleLength   = 7
	MaxTitleLength   = 77
	MinContentLength = 21
	MaxContentLength = 777
	maxUploadSize    = 20 * 1024 * 1024 // 20MB
)

func CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get userID from context
		currentUserID := context.MustGetUser(r.Context()).ID

		postErrors := &PostCreationErrors{}

		post := model.Post{
			Id:         uuid.New().String(),
			UserId:     currentUserID,
			Title:      r.FormValue("title"),
			Content:    r.FormValue("content"),
			Visibility: r.FormValue("postPrivacy"),
		}

		allowedFollowersJSON := r.FormValue("allowedFollowers")
		if allowedFollowersJSON != "" {
			err := json.Unmarshal([]byte(allowedFollowersJSON), &post.AllowedFollowers)
			if err != nil {
				http.Error(w, "Invalid allowedFollowers format", http.StatusBadRequest)
				return
			}
		}

		// handle file upload
		file, header, err := r.FormFile("postImage")
		if err == nil {
			defer file.Close()

			// validate file size
			if header.Size > maxUploadSize {
				postErrors.PostImage = "File too large (max 20MB)"
				http.Error(w, "File too large (max 20MB)", http.StatusBadRequest)
				return
			}

			// read the file bytes
			buff := make([]byte, 512)
			if _, err := file.Read(buff); err != nil {
				http.Error(w, "Invalid file", http.StatusBadRequest)
				return
			}

			// return the reader to the begining of the file
			if _, err := file.Seek(0, 0); err != nil {
				http.Error(w, "File error", http.StatusInternalServerError)
				return
			}

			// validate file type
			filetype := http.DetectContentType(buff)
			if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" {
				postErrors.PostImage = "Only JPEG, PNG and GIF images are allowed"
				http.Error(w, "Only JPEG, PNG and GIF images are allowed", http.StatusBadRequest)
				return
			}

			// generate a unique filename
			ext := filepath.Ext(header.Filename)
			filename := uuid.New().String() + ext

			// define where to save the file (create an "uploads" directory first)
			filePath := filepath.Join("../frontend/public/uploads", "posts", filename)
			os.MkdirAll(filepath.Dir(filePath), os.ModePerm)

			// create the file
			dst, err := os.Create(filePath)
			if err != nil {
				log.Println("error creating file destination for post image: ", err)
				http.Error(w, "An error occured, please try again later: ", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// copy the uploaded file to the destination
			if _, err := io.Copy(dst, file); err != nil {
				log.Println("failed to save post image to the destination file: ", err)
				http.Error(w, "An error occured, please try again later: ", http.StatusInternalServerError)
				return
			}

			relativePath := strings.TrimPrefix(filePath, "../frontend/public")
			post.ImageUrl = sql.NullString{String: relativePath, Valid: true}
		}

		postErrors, hasErrors := validatePost(post)
		if hasErrors {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(postErrors)
			return
		}

		// insert the main post
		_, postInsertErr := db.Exec(`
            INSERT INTO posts (id, user_id, title, content, visibility, post_image, created_at)
            VALUES (?, ?, ?, ?, ?, ?, ?)`,
			post.Id, post.UserId, post.Title, post.Content, post.Visibility, post.ImageUrl, post.CreatedAt)
		if postInsertErr != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		// handle private posts (visibility = "private")
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
	if visibility != "public" && visibility != "private" && visibility != "almostprivate" {
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
