package service

import (
	"backend/internal/model"
	"backend/internal/repository"
)

type GroupService struct {
	Repo *repository.GroupRepository
}

// NewGroupService creates and returns a new instance of GroupService.
func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{Repo: repo}
}

func (s *GroupService) CreateGroup(title, description, privacySetting string, creatorID string) (*model.Group, error) {
	// Start a transaction within the service layer
	tx, err := s.Repo.DB.Begin() // Access DB from repository
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Create the Group in the repository
	newGroup := &model.Group{
		Title:          title,
		Description:    description,
		CreatorID:      creatorID,
		PrivacySetting: privacySetting,
	}
	groupID, err := s.Repo.InsertGroup(tx, newGroup) // Pass transaction to repo
	if err != nil {
		return nil, err
	}
	newGroup.ID = groupID // Update the group ID after insertion

	// Add creator as group admin in the repository
	groupMember := &model.GroupMember{
		GroupID: groupID,
		UserID:  creatorID,
		Role:    "admin",
		Status:  "active",
	}
	err = s.Repo.InsertGroupMember(tx, groupMember) // Pass transaction to repo
	if err != nil {
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return newGroup, nil
}
