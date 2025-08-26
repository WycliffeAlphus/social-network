package model

import "time"

type Notification struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	ActorID    int       `json:"actor_id"`
	Type       string    `json:"type"`
	ContentID  int       `json:"content_id,omitempty"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}
