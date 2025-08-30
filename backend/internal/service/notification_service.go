package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
	"log"
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
func (s *NotificationService) CreateFollowRequestNotification(actorID, targetUserID string) error {
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


// CreateFollowAcceptedNotification creates a notification for a new follow request.
func (s *NotificationService) CreateFollowAcceptedNotification(actorID, targetUserID string) error {
	actor, err := s.userRepo.GetUserByID(actorID)
	if err != nil {
		return err
	}

	notification := &model.Notification{
		UserID:  targetUserID,
		ActorID: actorID,
		Type:    "follow_accepted",
		Message: fmt.Sprintf("%s %s accepted your follow request.", actor.FirstName, actor.LastName),
	}

	return s.repo.Create(notification)
}

// CreateGroupInviteNotification creates a notification for a group invitation.
func (s *NotificationService) CreateGroupInviteNotification(actorID, targetUserID string, groupID int) error {
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
func (s *NotificationService) GetByUserID(userID string) ([]*model.Notification, error) {
	return s.repo.GetByUserID(userID)
}

// MarkAllAsRead marks all of a user's notifications as read.
func (s *NotificationService) MarkAllAsRead(userID string) error {
	return s.repo.MarkAllAsRead(userID)
}

// CreateGroupJoinRequestNotification creates a notification for the group owner when a user requests to join.
func (s *NotificationService) CreateGroupJoinRequestNotification(actorID, groupOwnerID string, groupID int) error {
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

// CreateGroupJoinAcceptedNotification creates a notification for a new follow request.
func (s *NotificationService) CreateGroupJoinAcceptedNotification(actorID, targetUserID string, groupID int) error {
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
		Type:      "group_join_accepted",
		ContentID: groupID,
		Message:   fmt.Sprintf("%s %s has accepted your request to join the group '%s'.", actor.FirstName, actor.LastName, group.Title),
	}

	return s.repo.Create(notification)
}

// CreateGroupEventNotification creates a notification for all group members when an event is created.
func (s *NotificationService) CreateGroupEventNotification(actorID string, groupID, eventID int) error {
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
		// Don't notify the user who created the event
		if memberIDStr == actorID {
			continue
		}

		notification := &model.Notification{
			UserID:    memberIDStr,
			ActorID:   actorID,
			Type:      "group_event_created",
			ContentID: eventID,
			Message:   message,
		}

		if err := s.repo.Create(notification); err != nil {
			// Log the error but continue trying to notify other members
			log.Printf("Failed to create event notification for user %s: %v", memberIDStr, err)
		}
	}

	return nil
}

// CreatePostNotification creates a notification for a new post.
func (s *NotificationService) CreatePostNotification(actorID, postID string, groupID *int) error {
	actor, err := s.userRepo.GetUserByID(actorID)
	if err != nil {
		return fmt.Errorf("failed to get actor: %w", err)
	}

	var message string
	var members []string

	if groupID != nil {
		// It's a group post
		group, err := s.groupRepo.FindGroupByID(uint(*groupID))
		if err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}
		message = fmt.Sprintf("%s %s posted in %s.", actor.FirstName, actor.LastName, group.Title)
		members, err = s.groupRepo.GetGroupMembers(uint(*groupID))
		if err != nil {
			return fmt.Errorf("failed to get group members: %w", err)
		}
	} else {
		// It's a public post, notify followers
		message = fmt.Sprintf("%s %s created a new post.", actor.FirstName, actor.LastName)
		followers, err := s.userRepo.GetFollowers(actorID)
		if err != nil {
			return fmt.Errorf("failed to get followers: %w", err)
		}
		for _, follower := range followers {
			members = append(members, follower.ID)
		}
	}

	for _, memberID := range members {
		if memberID == actorID {
			continue // Don't notify the actor
		}

		notification := &model.Notification{
			UserID:  memberID,
			ActorID: actorID,
			Type:    "new_post",
			Message: message,
			// Assuming you might want to link to the post, you'd need a way to represent this.
			// Maybe ContentID could store postID, but it's an int. This needs schema adjustment.
		}

		if err := s.repo.Create(notification); err != nil {
			log.Printf("Failed to create post notification for user %s: %v", memberID, err)
		}
	}

	return nil
}

// CreateCommentNotification creates a notification for a new comment on a post.
func (s *NotificationService) CreateCommentNotification(actorID, postOwnerID, postID string) error {
	// Prevent self-notification
	if actorID == postOwnerID {
		return nil
	}

	actor, err := s.userRepo.GetUserByID(actorID)
	if err != nil {
		return fmt.Errorf("failed to get actor: %w", err)
	}

	message := fmt.Sprintf("%s %s commented on your post.", actor.FirstName, actor.LastName)

	notification := &model.Notification{
		UserID:  postOwnerID,
		ActorID: actorID,
		Type:    "new_comment",
		Message: message,
		// Again, linking to the post/comment would be ideal but requires schema changes.
	}

	return s.repo.Create(notification)
}

// CreateReactionNotification creates a notification for a new reaction on a post.
func (s *NotificationService) CreateReactionNotification(actorID, postOwnerID, postID string) error {
	// Prevent self-notification
	if actorID == postOwnerID {
		return nil
	}

	actor, err := s.userRepo.GetUserByID(actorID)
	if err != nil {
		return fmt.Errorf("failed to get actor: %w", err)
	}

	message := fmt.Sprintf("%s %s reacted to your post.", actor.FirstName, actor.LastName)

	notification := &model.Notification{
		UserID:  postOwnerID,
		ActorID: actorID,
		Type:    "new_reaction",
		Message: message,
	}

	return s.repo.Create(notification)
}
