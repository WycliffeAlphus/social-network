package model

import "time"

type Post struct {
	Id               string    `json:"id,omitempty"`
	UserId           string    `json:"userid,omitempty"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	Visibility       string    `json:"status"`
	ImageUrl         string    `json:"imageurl,omitempty"`
	CreatedAt        time.Time `json:"createdat"`
	AllowedFollowers []string  `json:"allowedfollowers,omitempty"`
}
