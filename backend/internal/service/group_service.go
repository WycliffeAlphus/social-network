package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
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
func (s *GroupService) RequestToJoinGroup(groupID uint, userID string) (*model.Group, error) {
	// Check if group exists
	group, err := s.Repo.FindGroupByID(groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	// Check if user is already a member or has a pending request
	isMember, status, err := s.Repo.CheckUserMembership(groupID, userID)
	if err != nil {
		return nil, err
	}
	if isMember {
		if status == "active" {
			return nil, fmt.Errorf("user is already a member of this group")
		} else if status == "pending" {
			return nil, fmt.Errorf("user already has a pending join request for this group")
		}
	}

	// Check if user is the group creator (creators are automatically members)
	if group.CreatorID == userID {
		return nil, fmt.Errorf("group creator cannot request to join their own group")
	}

	// Create the join request
	err = s.Repo.CreateJoinRequest(groupID, userID)
	if err != nil {
		return nil, err
	}

	return group, nil
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

// InviteUserToGroup creates an invitation for a user to join a group.
func (s *GroupService) InviteUserToGroup(groupID uint, inviterID string, targetUserID string) error {
	// In a real app, you would check if the inviterID has permission to invite.
	// For now, we'll assume any member can invite.

	// Check if target user is already a member
	isMember, _, err := s.Repo.CheckUserMembership(groupID, targetUserID)
	if err != nil {
		return err
	}
	if isMember {
		return fmt.Errorf("user is already a member or has a pending invitation")
	}

	// Create the invitation (pending member)
	return s.Repo.CreateGroupInvitation(groupID, inviterID, targetUserID)
}

// AcceptGroupInvitation accepts a group invitation.
func (s *GroupService) AcceptGroupInvitation(invitationID int, userID string) error {
	invite, err := s.Repo.GetGroupInvitation(invitationID)
	if err != nil {
		return err
	}
	if invite.InvitedID != userID {
		return fmt.Errorf("you are not authorized to accept this invitation")
	}
	if invite.Status != "pending" {
		return fmt.Errorf("this invitation is no longer valid")
	}

	tx, err := s.Repo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.Repo.AcceptGroupInvitation(invitationID); err != nil {
		return err
	}

	member := &model.GroupMember{
		GroupID: uint(invite.GroupID),
		UserID:  invite.InvitedID,
		Role:    "member",
		Status:  "active",
	}
	if err := s.Repo.InsertGroupMember(tx, member); err != nil {
		return err
	}

	return tx.Commit()
}

// DeclineGroupInvitation declines a group invitation.
func (s *GroupService) DeclineGroupInvitation(invitationID int, userID string) error {
	invite, err := s.Repo.GetGroupInvitation(invitationID)
	if err != nil {
		return err
	}
	if invite.InvitedID != userID {
		return fmt.Errorf("you are not authorized to decline this invitation")
	}
	if invite.Status != "pending" {
		return fmt.Errorf("this invitation is no longer valid")
	}

	return s.Repo.DeclineGroupInvitation(invitationID)
}

// CreateEvent creates a new event for a group.
func (s *GroupService) CreateEvent(event *model.Event) (*model.Event, error) {
	// In a real app, you would check if the creator is a member of the group.
	id, err := s.Repo.CreateEvent(event)
	if err != nil {
		return nil, err
	}
	event.ID = id
	return event, nil
}
