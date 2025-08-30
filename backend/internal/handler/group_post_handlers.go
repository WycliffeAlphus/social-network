package handler

import (
	"backend/internal/model"
	"backend/pkg/getusers"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract group ID from URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}
		groupID := pathParts[3]

		row := db.QueryRow("SELECT id, title, description FROM groups WHERE id = ?", groupID)

		var group model.Group
		err := row.Scan(&group.ID, &group.Title, &group.Description)
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "Failed to fetch group", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(group)
	}
}

func GetGroupPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract group ID from URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}
		groupID := pathParts[3]

		rows, err := db.Query(`
			SELECT
				p.id, p.user_id, p.title, p.content, p.visibility, p.post_image, p.created_at,
				(
					SELECT COUNT(1) FROM comments c WHERE c.post_id = p.id AND c.parent_id IS NULL
				) AS comment_count
			FROM posts p
			WHERE p.group_id = ?
			ORDER BY p.created_at DESC`, groupID)
		if err != nil {
			http.Error(w, "Failed to fetch group posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []model.Post
		for rows.Next() {
			var post model.Post
			if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.Visibility, &post.ImageUrl, &post.CreatedAt, &post.CommentCount); err != nil {
				http.Error(w, "Failed to scan post", http.StatusInternalServerError)
				return
			}

			// Populate Creator and CreatorImg
			user, err := getusers.GetUserByID(db, post.UserId)
			if err != nil {
				// Handle error, e.g., log it or skip this post
				fmt.Println("User not found for post:", post.Id, err)
				continue // Skip this post if user not found
			}
			post.Creator = user.FirstName + " " + user.LastName
			post.CreatorImg = user.ImgURL

			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
