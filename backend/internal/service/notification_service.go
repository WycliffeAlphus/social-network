package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"log"
	"strconv"
)

type NotificationService struct {
	repo      *repository.NotificationRepository
	userRepo  *repository.UserRepository
	groupRepo *repository.GroupRepository
}

func NewNotificationService(repo *repository.NotificationRepository, userRepo *repository.UserRepository, groupRepo *repository.GroupRepository) *NotificationService {
	return &NotificationService{repo: repo, userRepo: userRepo, groupRepo: groupRepo}
}

// CreateFollowRequestNotification creates a notification for a new follow request.
func (s *NotificationService) CreateFollowRequestNotification(actorID, targetUserID int) error {
	actor, err := s.userRepo.GetUserByID(actorID)
	if err != nil {
		return err
	}

	notification := &model.Notification{
		UserID:  targetUserID,
		ActorID: actorID,
		Type:    "follow_request",
		Message: fmt.Sprintf("%s %s wants to follow you.", actor.FirstName, actor.LastName),
	}

	return s.repo.Create(notification)
}

// CreateGroupInviteNotification creates a notification for a group invitation.
func (s *NotificationService) CreateGroupInviteNotification(actorID, targetUserID, groupID int) error {
	actor, err := s.userRepo.GetUserByID(actorID)
	if err != nil {
		return err
	}
	group, err := s.groupRepo.FindGroupByID(uint(groupID))
	if err != nil {
		return err
	}

	notification := &model.Notification{
		UserID:    targetUserID,
		ActorID:   actorID,
		Type:      "group_invite",
		ContentID: groupID,
		Message:   fmt.Sprintf("%s %s has invited you to join the group '%s'.", actor.FirstName, actor.LastName, group.Title),
	}

	return s.repo.Create(notification)
}

// GetByUserID retrieves all notifications for a specific user.
func (s *NotificationService) GetByUserID(userID int) ([]*model.Notification, error) {
	return s.repo.GetByUserID(userID)
}

// MarkAllAsRead marks all of a user's notifications as read.
func (s *NotificationService) MarkAllAsRead(userID int) error {
	return s.repo.MarkAllAsRead(userID)
}

// CreateGroupJoinRequestNotification creates a notification for the group owner when a user requests to join.
func (s *NotificationService) CreateGroupJoinRequestNotification(actorID, groupOwnerID, groupID int) error {
	actor, err := s.userRepo.GetUserByID(actorID)
	if err != nil {
		return err
	}
	group, err := s.groupRepo.FindGroupByID(uint(groupID))
	if err != nil {
		return err
	}

	notification := &model.Notification{
		UserID:    groupOwnerID,
		ActorID:   actorID,
		Type:      "group_join_request",
		ContentID: groupID,
		Message:   fmt.Sprintf("%s %s has requested to join your group '%s'.", actor.FirstName, actor.LastName, group.Title),
	}

	return s.repo.Create(notification)
}

// CreateGroupEventNotification creates a notification for all group members when an event is created.
func (s *NotificationService) CreateGroupEventNotification(actorID, groupID, eventID int) error {
	actor, err := s.userRepo.GetUserByID(actorID)
	if err != nil {
		return fmt.Errorf("failed to get actor: %w", err)
	}
	group, err := s.groupRepo.FindGroupByID(uint(groupID))
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	members, err := s.groupRepo.GetGroupMembers(uint(groupID))
	if err != nil {
		return fmt.Errorf("failed to get group members: %w", err)
	}

	// In a real-world scenario, you would fetch the event title as well.
	// For now, we'll use a generic message.
	message := fmt.Sprintf("%s %s has created a new event in '%s'.", actor.FirstName, actor.LastName, group.Title)

	for _, memberIDStr := range members {
		memberID, _ := strconv.Atoi(memberIDStr)
		// Don't notify the user who created the event
		if memberID == actorID {
			continue
		}

		notification := &model.Notification{
			UserID:    memberID,
			ActorID:   actorID,
			Type:      "group_event_created",
			ContentID: eventID,
			Message:   message,
		}

		if err := s.repo.Create(notification); err != nil {
			// Log the error but continue trying to notify other members
			log.Printf("Failed to create event notification for user %d: %v", memberID, err)
		}
	}

	return nil
}
