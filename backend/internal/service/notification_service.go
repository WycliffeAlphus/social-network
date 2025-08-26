package service

import (
	"fmt"
	"social-network/internal/model"
	"social-network/internal/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
	userRepo *repository.UserRepository // Assuming you have a user repository to get user names
}

func NewNotificationService(repo *repository.NotificationRepository, userRepo *repository.UserRepository) *NotificationService {
	return &NotificationService{repo: repo, userRepo: userRepo}
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
	// Implementation to be added: Get actor name and group name
	notification := &model.Notification{
		UserID:  targetUserID,
		ActorID: actorID,
		Type:    "group_invite",
		ContentID: groupID,
		Message: fmt.Sprintf("You have been invited to join a group."), // Placeholder message
	}

	return s.repo.Create(notification)
}
