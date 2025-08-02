package model

import "time"

type GroupMember struct {
	ID        uint       `json:"id"`
	GroupID   uint       `json:"group_id"`
	UserID    string     `json:"user_id"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	JoinedAt  time.Time  `json:"joined_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
