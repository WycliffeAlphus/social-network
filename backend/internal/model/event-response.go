package model

import "time"

type EventResponse struct {
	ID        uint       `json:"id"`
	EventID   uint       `json:"event_id"`
	UserID    string     `json:"user_id"`
	Response  string     `json:"response"` // 'going', 'not_going'
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type EventWithResponses struct {
	Event
	GoingCount    int `json:"going_count"`
	NotGoingCount int `json:"not_going_count"`
	UserResponse  string `json:"user_response,omitempty"`
}
