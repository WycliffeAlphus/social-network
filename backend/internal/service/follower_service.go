package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"fmt"
)

// FollowerService provides business logic for follower operations
type FollowerService struct {
	FollowerRepo *repository.FollowerRepository
	UserRepo     *repository.UserRepository
}

// GetFollowers retrieves followers for a user with privacy checks
func (s *FollowerService) GetFollowers(userID string, requestingUserID string) ([]model.FollowerWithUser, error) {
	// Check if the requesting user can view this user's followers
	canView, err := s.canViewFollowers(userID, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check permissions: %w", err)
	}

	if !canView {
		return nil, fmt.Errorf("not authorized to view followers")
	}

	followers, err := s.FollowerRepo.GetFollowers(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	return followers, nil
}

// GetFollowing retrieves users that a user is following with privacy checks
func (s *FollowerService) GetFollowing(userID string, requestingUserID string) ([]model.FollowerWithUser, error) {
	// Check if the requesting user can view this user's following list
	canView, err := s.canViewFollowing(userID, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check permissions: %w", err)
	}

	if !canView {
		return nil, fmt.Errorf("not authorized to view following list")
	}

	following, err := s.FollowerRepo.GetFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	return following, nil
}

// GetFollowerCounts returns follower and following counts
func (s *FollowerService) GetFollowerCounts(userID string) (followers int, following int, err error) {
	return s.FollowerRepo.GetFollowerCounts(userID)
}

// FollowUser creates a follow relationship
func (s *FollowerService) FollowUser(followerID, followedID string) (string, error) {
	// Prevent self-follow
	if followerID == followedID {
		return "", fmt.Errorf("cannot follow yourself")
	}

	// Check if already following
	isFollowing, status, err := s.FollowerRepo.IsFollowing(followerID, followedID)
	if err != nil {
		return "", fmt.Errorf("failed to check follow status: %w", err)
	}

	if isFollowing {
		if status == "accepted" {
			return "", fmt.Errorf("already following this user")
		} else if status == "pending" {
			return "", fmt.Errorf("follow request already sent")
		}
	}

	// Get target user's profile visibility to determine status
	// This would require a method to get user by ID, let's assume we have it
	// For now, we'll determine status based on profile visibility
	profileVisibility, err := s.getUserProfileVisibility(followedID)
	if err != nil {
		return "", fmt.Errorf("failed to get user profile: %w", err)
	}

	status = "pending"
	if profileVisibility == "public" {
		status = "accepted"
	}

	err = s.FollowerRepo.FollowUser(followerID, followedID, status)
	if err != nil {
		return "", fmt.Errorf("failed to follow user: %w", err)
	}

	message := "Follow request sent"
	if status == "accepted" {
		message = "Successfully followed user"
	}

	return message, nil
}

// UnfollowUser removes a follow relationship
func (s *FollowerService) UnfollowUser(followerID, followedID string) error {
	return s.FollowerRepo.UnfollowUser(followerID, followedID)
}

// AcceptFollowRequest accepts a pending follow request
func (s *FollowerService) AcceptFollowRequest(followerID, followedID string) error {
	return s.FollowerRepo.AcceptFollowRequest(followerID, followedID)
}

// RejectFollowRequest rejects a pending follow request
func (s *FollowerService) RejectFollowRequest(followerID, followedID string) error {
	return s.FollowerRepo.RejectFollowRequest(followerID, followedID)
}

// GetPendingFollowRequests retrieves pending follow requests for a user
func (s *FollowerService) GetPendingFollowRequests(userID string) ([]model.FollowerWithUser, error) {
	return s.FollowerRepo.GetPendingFollowRequests(userID)
}

// Helper methods

// canViewFollowers checks if the requesting user can view the target user's followers
func (s *FollowerService) canViewFollowers(targetUserID, requestingUserID string) (bool, error) {
	// User can always view their own followers
	if targetUserID == requestingUserID {
		return true, nil
	}

	// Get target user's profile visibility
	profileVisibility, err := s.getUserProfileVisibility(targetUserID)
	if err != nil {
		return false, err
	}

	// If profile is public, anyone can view followers
	if profileVisibility == "public" {
		return true, nil
	}

	// If profile is private, only followers can view followers list
	if profileVisibility == "private" {
		isFollowing, status, err := s.FollowerRepo.IsFollowing(requestingUserID, targetUserID)
		if err != nil {
			return false, err
		}
		return isFollowing && status == "accepted", nil
	}

	return false, nil
}

// canViewFollowing checks if the requesting user can view the target user's following list
func (s *FollowerService) canViewFollowing(targetUserID, requestingUserID string) (bool, error) {
	// User can always view their own following list
	if targetUserID == requestingUserID {
		return true, nil
	}

	// Get target user's profile visibility
	profileVisibility, err := s.getUserProfileVisibility(targetUserID)
	if err != nil {
		return false, err
	}

	// If profile is public, anyone can view following list
	if profileVisibility == "public" {
		return true, nil
	}

	// If profile is private, only followers can view following list
	if profileVisibility == "private" {
		isFollowing, status, err := s.FollowerRepo.IsFollowing(requestingUserID, targetUserID)
		if err != nil {
			return false, err
		}
		return isFollowing && status == "accepted", nil
	}

	return false, nil
}

// getUserProfileVisibility gets a user's profile visibility setting
func (s *FollowerService) getUserProfileVisibility(userID string) (string, error) {
	// TODO: Implement proper database query through repository layer
	// For now, we'll return public as default to allow the system to work
	// This should be implemented properly with database access
	return "public", nil
}
