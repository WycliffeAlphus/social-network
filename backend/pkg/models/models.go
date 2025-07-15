package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                string
	Email             string
	Password          string
	FirstName         string
	LastName          string
	DateOfBirth       time.Time
	AvatarImage       sql.NullString
	Nickname          sql.NullString
	AboutMe           sql.NullString
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ProfileVisibility string
}
