package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) RegisterUser(user *model.User) error {
	user.ID = utils.GenerateUUID()
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.Repo.CreateUser(user)
}
