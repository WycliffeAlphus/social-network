package model

type Message struct {
	Type      string `json:"type,omitempty"` // "message" or "typing"
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	IsTyping  bool   `json:"isTyping,omitempty"`
}