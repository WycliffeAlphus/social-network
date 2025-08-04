package model

import (
	"database/sql"
	"time"
)

type Post struct {
	Id               string         `json:"id,omitempty"`
	UserId           string         `json:"userid,omitempty"`
	Title            string         `json:"title"`
	Content          string         `json:"content"`
	Visibility       string         `json:"status"`
	ImageUrl         sql.NullString `json:"imageurl"`
	CreatedAt        time.Time      `json:"createdat"`
	AllowedFollowers []string       `json:"allowedfollowers,omitempty"`
}
