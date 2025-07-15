package model

import "time"

type User struct {
	ID                string    `json:"id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	DOB               time.Time `json:"dob"`
	ImgURL            string    `json:"img_url,omitempty"`
	Nickname          string    `json:"nickname,omitempty"`
	About             string    `json:"about,omitempty"`
	Password          string    `json:"password"`
	ProfileVisibility string    `json:"profile_visibility"`
	CreatedAt         time.Time `json:"created_at"`
}
