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

// GetUserGroups retrieves all groups where the user is an active member.
func (s *GroupService) GetUserGroups(userID string) ([]model.Group, error) {
	return s.Repo.GetUserGroups(userID)
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
		if status == "active" {
			return fmt.Errorf("user is already a member of this group")
		} else if status == "pending" {
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

// GetPendingJoinRequests retrieves all pending join requests for groups created by the user.
func (s *GroupService) GetPendingJoinRequests(groupID uint, creatorUserID string) ([]model.GroupJoinRequest, error) {
	// Verify that the user requesting is the group creator
	isCreator, err := s.Repo.IsGroupCreator(groupID, creatorUserID)
	if err != nil {
		fmt.Printf("[DEBUG] GetPendingJoinRequests - IsGroupCreator error: %v\n", err)
		return nil, err
	}
	fmt.Printf("[DEBUG] GetPendingJoinRequests - isCreator: %v for groupID: %d, userID: %s\n", isCreator, groupID, creatorUserID)
	if !isCreator {
		return nil, fmt.Errorf("only group creators can view join requests")
	}

	// Get pending requests for this group
	requests, err := s.Repo.GetPendingJoinRequests(groupID)
	fmt.Printf("[DEBUG] GetPendingJoinRequests - repo returned %d requests\n", len(requests))
	return requests, err
}

// RejectJoinRequest allows a group creator to reject a pending join request.
func (s *GroupService) RejectJoinRequest(groupID uint, requesterUserID string, creatorUserID string) error {
	// Verify that the user rejecting the request is the group creator
	isCreator, err := s.Repo.IsGroupCreator(groupID, creatorUserID)
	if err != nil {
		return err
	}
	if !isCreator {
		return fmt.Errorf("only group creators can reject join requests")
	}

	// Check if there's a pending request for this user
	isMember, status, err := s.Repo.CheckUserMembership(groupID, requesterUserID)
	if err != nil {
		return err
	}
	if !isMember || status != "pending" {
		return fmt.Errorf("no pending join request found for this user")
	}

	// Reject the join request
	return s.Repo.RejectJoinRequest(groupID, requesterUserID)
}

// IsUserGroupMember checks if a user is an active member of a group.
func (s *GroupService) IsUserGroupMember(groupID uint, userID string) (bool, error) {
	isMember, status, err := s.Repo.CheckUserMembership(groupID, userID)
	if err != nil {
		return false, err
	}
	return isMember && status == "active", nil
}

// InviteUserToGroup allows a group creator or member to invite a user to join a group.
func (s *GroupService) InviteUserToGroup(groupID uint, invitedUserID string, inviterUserID string) error {
	// Verify that the user sending the invite is a member of the group
	isMember, status, err := s.Repo.CheckUserMembership(groupID, inviterUserID)
	if err != nil {
		return err
	}
	if !isMember || status != "active" {
		return fmt.Errorf("only group members can invite users")
	}

	// Check if group exists
	group, err := s.Repo.FindGroupByID(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	// Check if user is already a member or has a pending request/invite
	isAlreadyMember, memberStatus, err := s.Repo.CheckUserMembership(groupID, invitedUserID)
	if err != nil {
		return err
	}
	if isAlreadyMember {
		if memberStatus == "active" {
			return fmt.Errorf("user is already a member of this group")
		} else if memberStatus == "pending" {
			return fmt.Errorf("user already has a pending request for this group")
		} else if memberStatus == "invited" {
			return fmt.Errorf("user has already been invited to this group")
		}
	}

	// Create the invite
	return s.Repo.CreateGroupInvite(groupID, invitedUserID)
}

// AcceptGroupInvite allows an invited user to accept an invitation to join a group.
func (s *GroupService) AcceptGroupInvite(groupID uint, userID string) error {
	// Check if there's a pending invitation for this user
	isMember, status, err := s.Repo.CheckUserMembership(groupID, userID)
	if err != nil {
		return err
	}
	if !isMember || status != "invited" {
		return fmt.Errorf("no pending invitation found for this user")
	}

	// Accept the invitation
	return s.Repo.AcceptGroupInvite(groupID, userID)
}

// RejectGroupInvite allows an invited user to reject an invitation to join a group.
func (s *GroupService) RejectGroupInvite(groupID uint, userID string) error {
	// Check if there's a pending invitation for this user
	isMember, status, err := s.Repo.CheckUserMembership(groupID, userID)
	if err != nil {
		return err
	}
	if !isMember || status != "invited" {
		return fmt.Errorf("no pending invitation found for this user")
	}

	// Reject the invitation
	return s.Repo.RejectGroupInvite(groupID, userID)
}

// GetUserGroupInvites retrieves all pending invitations for a user.
func (s *GroupService) GetUserGroupInvites(userID string) ([]model.Group, error) {
	return s.Repo.GetUserGroupInvites(userID)
}
