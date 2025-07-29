package model

import "time"

type Group struct {
	ID             uint       `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	CreatorID      uint       `json:"creator_id"`
	PrivacySetting string     `json:"privacy_setting"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
