package model

import (
	"database/sql"
	"time"
)

type Comment struct {
	Id        string         `json:"id"`
	PostId    string         `json:"post_id"`
	UserId    string         `json:"user_id"`
	Content   string         `json:"content"`
	ParentId  sql.NullString `json:"parent_id,omitempty"` // For replies to other comments
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Additional fields for API responses
	UserFirstName string    `json:"user_first_name,omitempty"`
	UserLastName  string    `json:"user_last_name,omitempty"`
	UserNickname  string    `json:"user_nickname,omitempty"`
	UserImgURL    string    `json:"user_img_url,omitempty"`
	Replies       []Comment `json:"replies,omitempty"`
}

// CommentRequest represents the request body for creating a comment
type CommentRequest struct {
	Content  string `json:"content"`
	ParentId string `json:"parent_id,omitempty"` // For replies
}

// CommentResponse represents the response for comment operations
type CommentResponse struct {
	Success bool                 `json:"success"`
	Comment *CommentWithUserInfo `json:"comment,omitempty"`
	Error   string               `json:"error,omitempty"`
}

// CommentWithUserInfo represents a comment with user information for API responses
type CommentWithUserInfo struct {
	Id        string    `json:"id"`
	PostId    string    `json:"post_id"`
	UserId    string    `json:"user_id"`
	Content   string    `json:"content"`
	ParentId  *string   `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// User info
	UserFirstName string `json:"user_first_name"`
	UserLastName  string `json:"user_last_name"`
	UserNickname  string `json:"user_nickname"`
	UserImgURL    string `json:"user_img_url"`

	// Nested replies
	Replies []CommentWithUserInfo `json:"replies,omitempty"`
}
