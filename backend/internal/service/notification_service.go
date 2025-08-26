package service

import (
	"fmt"
	"social-network/internal/model"
	"social-network/internal/repository"
)

type NotificationService struct {
	repo     *repository.NotificationRepository
	userRepo *repository.UserRepository
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
	// Implementation to be added
	return nil
}

// CreateGroupEventNotification creates a notification for all group members when an event is created.
func (s *NotificationService) CreateGroupEventNotification(actorID, groupID, eventID int) error {
	// Implementation to be added
	return nil
}
