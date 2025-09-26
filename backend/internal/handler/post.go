package handler

import (
	"backend/internal/context"
	"database/sql"
	"encoding/json"
	"net/http"

	"backend/internal/repository"
	"backend/pkg/extractid"
)

func GetPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		postID := extractid.ExtractUserIDFromPath(r.URL.Path, "post")
		if postID == "" {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		userID := context.MustGetUser(r.Context()).ID

		post, err := repository.GetPostByID(db, postID, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}
