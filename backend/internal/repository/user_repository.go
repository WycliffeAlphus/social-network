package repository

import (
	"backend/internal/model"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CreateUser(user *model.User) error {
	  query := `INSERT INTO users (id, email, password, first_name, last_name, dob, avatar, nickname, about_me)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
			    _, err := r.DB.Exec(query, user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.DOB, user.Avatar, user.Nickname, user.AboutMe)
    return err

}
