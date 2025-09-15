package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"time"
)

type GroupService struct {
	Repo *repository.GroupRepository
}

// NewGroupService creates and returns a new instance of GroupService.
func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{Repo: repo}
}

// GetAllGroups retrieves all groups from the repository.
func (s *GroupService) GetAllGroups() ([]model.Group, error) {
	return s.Repo.FindAll()
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

// RequestToJoinGroup creates a join request for a user to join a group.
func (s *GroupService) RequestToJoinGroup(groupID uint, userID string) error {
	// Check if group exists
	group, err := s.Repo.FindGroupByID(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	// Check if user is already a member or has a pending request
	isMember, status, err := s.Repo.CheckUserMembership(groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		switch status {
		case "active":
			return fmt.Errorf("user is already a member of this group")
		case "pending":
			return fmt.Errorf("user already has a pending join request for this group")
		}
	}

	// Check if user is the group creator (creators are automatically members)
	if group.CreatorID == userID {
		return fmt.Errorf("group creator cannot request to join their own group")
	}

	// Create the join request
	return s.Repo.CreateJoinRequest(groupID, userID)
}

// AcceptJoinRequest allows a group creator to accept a pending join request.
func (s *GroupService) AcceptJoinRequest(groupID uint, requesterUserID string, creatorUserID string) error {
	// Verify that the user accepting the request is the group creator
	isCreator, err := s.Repo.IsGroupCreator(groupID, creatorUserID)
	if err != nil {
		return err
	}
	if !isCreator {
		return fmt.Errorf("only group creators can accept join requests")
	}

	// Check if there's a pending request for this user
	isMember, status, err := s.Repo.CheckUserMembership(groupID, requesterUserID)
	if err != nil {
		return err
	}
	if !isMember || status != "pending" {
		return fmt.Errorf("no pending join request found for this user")
	}

	// Accept the join request
	return s.Repo.AcceptJoinRequest(groupID, requesterUserID)
}

// CreateGroupEvent validates membership and creates a new event for the group
func (s *GroupService) CreateGroupEvent(groupID uint, creatorUserID string, title string, description string, t time.Time, location string) (*model.GroupEvent, error) {
	// ensure group exists
	group, err := s.Repo.FindGroupByID(groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	// only active members can create events (including creator)
	isMember, err := s.Repo.IsActiveMember(groupID, creatorUserID)
	if err != nil {
		return nil, err
	}
	if !isMember && group.CreatorID != creatorUserID {
		return nil, fmt.Errorf("only active members can create events")
	}

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	event := &model.GroupEvent{
		GroupID:     groupID,
		Title:       title,
		Description: description,
		Time:        t,
		Location:    location,
	}
	id, err := s.Repo.InsertGroupEvent(event)
	if err != nil {
		return nil, err
	}
	event.ID = id
	return event, nil
}

func (s *GroupService) ViewGroup(groupid uint) (model.Group, error) {
	groupDetails, err := s.Repo.FindGroupByID(groupid)
	return *groupDetails, err
}
