package model

import "time"

type GroupInvite struct {
	ID            uint       `json:"id"`
	GroupID       uint       `json:"group_id"`
	InviterUserID uint       `json:"inviter_user_id"`
	InvitedUserID uint       `json:"invited_user_id"`
	Status        string     `json:"status"`
	InvitedAt     time.Time  `json:"invited_at"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}
