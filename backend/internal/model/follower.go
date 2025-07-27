package model

import "time"

// FollowerWithUser represents a follower relationship with user details
type FollowerWithUser struct {
	FollowerID        string    `json:"follower_id"`
	FollowedID        string    `json:"followed_id"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UserID            string    `json:"user_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	Nickname          string    `json:"nickname,omitempty"`
	ImgURL            string    `json:"img_url,omitempty"`
	ProfileVisibility string    `json:"profile_visibility"`
}

// FollowersResponse represents the response for followers/following lists
type FollowersResponse struct {
	Status string              `json:"status"`
	Data   []FollowerWithUser  `json:"data"`
	Count  int                 `json:"count"`
}
