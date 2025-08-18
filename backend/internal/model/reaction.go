package model


type Reaction struct {
	PostID string  `json:"post_id"`
	UserID string  `json:"user_id"`
	Type   string `json:"type"` // "like" or "dislike"
}