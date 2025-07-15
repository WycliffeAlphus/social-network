package model

import "time"

type User struct {
	ID                string    `json:"id" db:"id"`
	Email             string    `json:"email" db:"email"`
	FirstName         string    `json:"first_name" db:"fname"`
	LastName          string    `json:"last_name" db:"lname"`
	DOB               time.Time `json:"dob" db:"dob"`
	ImgURL            string    `json:"img_url,omitempty" db:"imgurl"`
	Nickname          string    `json:"nickname,omitempty" db:"nickname"`
	About             string    `json:"about,omitempty" db:"about"`
	Password          string    `json:"password" db:"password"`
	ProfileVisibility string    `json:"profile_visibility" db:"profileVisibility"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}
