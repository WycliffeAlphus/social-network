package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MinCommentLength = 1
	MaxCommentLength = 500
)

// CommentHandler handles both GET and POST requests for /posts/:id/comments
func CommentHandler(db *sql.DB, notificationService *service.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the path ends with /comments
		if !strings.HasSuffix(r.URL.Path, "/comments") {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodPost:
			createComment(w, r, db, notificationService)
		case http.MethodGet:
			getComments(w, r, db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// createComment handles POST /posts/:id/comments
func createComment(w http.ResponseWriter, r *http.Request, db *sql.DB, notificationService *service.NotificationService) {
	// Get user ID from context
	currentUser := context.MustGetUser(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract post ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	postId := pathParts[3] // /posts/:id/comments

	// Parse request body
	var req model.CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate comment content
	if err := validateComment(req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if post exists
	var postExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postId).Scan(&postExists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !postExists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// If this is a reply, validate parent comment exists
	if req.ParentId != "" {
		var parentExists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM comments WHERE id = ? AND post_id = ?)", req.ParentId, postId).Scan(&parentExists)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if !parentExists {
			http.Error(w, "Parent comment not found", http.StatusNotFound)
			return
		}
	}

	// Create comment
	commentId := uuid.New().String()
	now := time.Now()

	var parentId sql.NullString
	if req.ParentId != "" {
		parentId = sql.NullString{String: req.ParentId, Valid: true}
	}

	_, err = db.Exec(`
		INSERT INTO comments (id, post_id, user_id, content, parent_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		commentId, postId, currentUser.ID, req.Content, parentId, now, now)

	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// Get the created comment with user info
	comment, err := getCommentWithUserInfo(db, commentId)
	if err != nil {
		http.Error(w, "Failed to retrieve comment", http.StatusInternalServerError)
		return
	}

	// Get post owner ID
	postOwnerID, err := repository.GetPostOwnerID(db, postId)
	if err != nil {
		log.Printf("Failed to get post owner ID: %v", err)
	} else {
		// Create notification
		if err := notificationService.CreateCommentNotification(currentUser.ID, postOwnerID, postId); err != nil {
			log.Printf("Failed to create comment notification: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.CommentResponse{
		Success: true,
		Comment: comment,
	})
}

// getComments handles GET /posts/:id/comments
func getComments(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract post ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	postId := pathParts[3] // /posts/:id/comments

	// Check if post exists
	var postExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postId).Scan(&postExists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !postExists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Get all comments for the post (top-level comments only)
	rows, err := db.Query(`
		SELECT c.id, c.post_id, c.user_id, c.content, c.parent_id, c.created_at, c.updated_at,
			   u.fname, u.lname, u.nickname, u.imgurl
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ? AND c.parent_id IS NULL
		ORDER BY c.created_at ASC`,
		postId)
	if err != nil {
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []model.CommentWithUserInfo
	for rows.Next() {
		var comment model.CommentWithUserInfo
		var parentId sql.NullString
		var userFirstName, userLastName, userNickname, userImgURL sql.NullString

		err := rows.Scan(
			&comment.Id, &comment.PostId, &comment.UserId, &comment.Content,
			&parentId, &comment.CreatedAt, &comment.UpdatedAt,
			&userFirstName, &userLastName, &userNickname, &userImgURL)
		if err != nil {
			continue
		}

		if parentId.Valid {
			comment.ParentId = &parentId.String
		}

		comment.UserFirstName = userFirstName.String
		comment.UserLastName = userLastName.String
		comment.UserNickname = userNickname.String
		comment.UserImgURL = userImgURL.String

		// Get replies for this comment
		commentReplies, err := buildRepliesTree(db, comment.Id)
		if err == nil {
			comment.Replies = commentReplies
		}

		comments = append(comments, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"comments": comments,
	})
}

// buildRepliesTree constructs a nested list of replies for a given parent comment id
func buildRepliesTree(db *sql.DB, parentCommentId string) ([]model.CommentWithUserInfo, error) {
	replies, err := getRepliesForComment(db, parentCommentId)
	if err != nil {
		return nil, err
	}

	var nestedReplies []model.CommentWithUserInfo
	for _, reply := range replies {
		var parentId *string
		if reply.ParentId.Valid {
			parentId = &reply.ParentId.String
		}

		// Recursively build children for this reply
		children, _ := buildRepliesTree(db, reply.Id)

		nestedReplies = append(nestedReplies, model.CommentWithUserInfo{
			Id:            reply.Id,
			PostId:        reply.PostId,
			UserId:        reply.UserId,
			Content:       reply.Content,
			ParentId:      parentId,
			CreatedAt:     reply.CreatedAt,
			UpdatedAt:     reply.UpdatedAt,
			UserFirstName: reply.UserFirstName,
			UserLastName:  reply.UserLastName,
			UserNickname:  reply.UserNickname,
			UserImgURL:    reply.UserImgURL,
			Replies:       children,
		})
	}

	return nestedReplies, nil
}

func validateComment(content string) error {
	if len(strings.TrimSpace(content)) < MinCommentLength {
		return fmt.Errorf("Comment must be at least %d character long", MinCommentLength)
	}
	if len(content) > MaxCommentLength {
		return fmt.Errorf("Comment must be at most %d characters long", MaxCommentLength)
	}
	return nil
}

func getCommentWithUserInfo(db *sql.DB, commentId string) (*model.CommentWithUserInfo, error) {
	var comment model.CommentWithUserInfo
	var parentId sql.NullString
	var userFirstName, userLastName, userNickname, userImgURL sql.NullString

	err := db.QueryRow(`
		SELECT c.id, c.post_id, c.user_id, c.content, c.parent_id, c.created_at, c.updated_at,
			   u.fname, u.lname, u.nickname, u.imgurl
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?`,
		commentId).Scan(
		&comment.Id, &comment.PostId, &comment.UserId, &comment.Content,
		&parentId, &comment.CreatedAt, &comment.UpdatedAt,
		&userFirstName, &userLastName, &userNickname, &userImgURL)

	if err != nil {
		return nil, err
	}

	if parentId.Valid {
		comment.ParentId = &parentId.String
	}

	comment.UserFirstName = userFirstName.String
	comment.UserLastName = userLastName.String
	comment.UserNickname = userNickname.String
	comment.UserImgURL = userImgURL.String

	return &comment, nil
}

func getRepliesForComment(db *sql.DB, commentId string) ([]model.Comment, error) {
	rows, err := db.Query(`
		SELECT c.id, c.post_id, c.user_id, c.content, c.parent_id, c.created_at, c.updated_at,
			   u.fname, u.lname, u.nickname, u.imgurl
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.parent_id = ?
		ORDER BY c.created_at ASC`,
		commentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []model.Comment
	for rows.Next() {
		var reply model.Comment
		var parentId sql.NullString
		var userFirstName, userLastName, userNickname, userImgURL sql.NullString

		err := rows.Scan(
			&reply.Id, &reply.PostId, &reply.UserId, &reply.Content,
			&parentId, &reply.CreatedAt, &reply.UpdatedAt,
			&userFirstName, &userLastName, &userNickname, &userImgURL)
		if err != nil {
			continue
		}

		if parentId.Valid {
			reply.ParentId = parentId
		}

		reply.UserFirstName = userFirstName.String
		reply.UserLastName = userLastName.String
		reply.UserNickname = userNickname.String
		reply.UserImgURL = userImgURL.String

		replies = append(replies, reply)
	}

	return replies, nil
}