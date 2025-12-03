package model

import "time"

type GroupJoinRequest struct {
	GroupID      uint      `json:"group_id"`
	UserID       string    `json:"user_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	UserName     string    `json:"user_name"`
	ImgURL       string    `json:"img_url"`
	UserImageURL string    `json:"user_image_url"`
	RequestedAt  time.Time `json:"requested_at"`
}
