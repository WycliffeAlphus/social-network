package repository

import (
	"backend/internal/model"
	"database/sql"
	"fmt"
)

// FollowerRepository handles database operations for followers
type FollowerRepository struct {
	DB *sql.DB
}

// GetFollowers retrieves all followers for a given user ID
func (r *FollowerRepository) GetFollowers(userID string) ([]model.FollowerWithUser, error) {
	query := `
		SELECT 
			f.follower_id,
			f.followed_id,
			f.status,
			f.created_at,
			u.id,
			u.fname,
			u.lname,
			u.email,
			COALESCE(u.nickname, '') as nickname,
			COALESCE(u.imgurl, '') as imgurl,
			u.profileVisibility
		FROM followers f
		JOIN users u ON f.follower_id = u.id
		WHERE f.followed_id = ? AND f.status = 'accepted'
		ORDER BY f.created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query followers: %w", err)
	}
	defer rows.Close()

	var followers []model.FollowerWithUser
	for rows.Next() {
		var follower model.FollowerWithUser
		err := rows.Scan(
			&follower.FollowerID,
			&follower.FollowedID,
			&follower.Status,
			&follower.CreatedAt,
			&follower.UserID,
			&follower.FirstName,
			&follower.LastName,
			&follower.Email,
			&follower.Nickname,
			&follower.ImgURL,
			&follower.ProfileVisibility,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follower: %w", err)
		}
		followers = append(followers, follower)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating followers: %w", err)
	}

	return followers, nil
}

// GetFollowing retrieves all users that a given user is following
func (r *FollowerRepository) GetFollowing(userID string) ([]model.FollowerWithUser, error) {
	query := `
		SELECT 
			f.follower_id,
			f.followed_id,
			f.status,
			f.created_at,
			u.id,
			u.fname,
			u.lname,
			u.email,
			COALESCE(u.nickname, '') as nickname,
			COALESCE(u.imgurl, '') as imgurl,
			u.profileVisibility
		FROM followers f
		JOIN users u ON f.followed_id = u.id
		WHERE f.follower_id = ? AND f.status = 'accepted'
		ORDER BY f.created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query following: %w", err)
	}
	defer rows.Close()

	var following []model.FollowerWithUser
	for rows.Next() {
		var follow model.FollowerWithUser
		err := rows.Scan(
			&follow.FollowerID,
			&follow.FollowedID,
			&follow.Status,
			&follow.CreatedAt,
			&follow.UserID,
			&follow.FirstName,
			&follow.LastName,
			&follow.Email,
			&follow.Nickname,
			&follow.ImgURL,
			&follow.ProfileVisibility,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan following: %w", err)
		}
		following = append(following, follow)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating following: %w", err)
	}

	return following, nil
}

// GetFollowerCounts returns the count of followers and following for a user
func (r *FollowerRepository) GetFollowerCounts(userID string) (followers int, following int, err error) {
	// Get followers count
	err = r.DB.QueryRow("SELECT COUNT(*) FROM followers WHERE followed_id = ? AND status = 'accepted'", userID).Scan(&followers)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get followers count: %w", err)
	}

	// Get following count
	err = r.DB.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND status = 'accepted'", userID).Scan(&following)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get following count: %w", err)
	}

	return followers, following, nil
}

// IsFollowing checks if userID is following targetUserID
func (r *FollowerRepository) IsFollowing(userID, targetUserID string) (bool, string, error) {
	var status string
	err := r.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", userID, targetUserID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to check following status: %w", err)
	}
	return true, status, nil
}

// FollowUser creates a follow relationship
func (r *FollowerRepository) FollowUser(followerID, followedID, status string) error {
	query := "INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)"
	_, err := r.DB.Exec(query, followerID, followedID, status)
	if err != nil {
		return fmt.Errorf("failed to create follow relationship: %w", err)
	}
	return nil
}

// UnfollowUser removes a follow relationship
func (r *FollowerRepository) UnfollowUser(followerID, followedID string) error {
	query := "DELETE FROM followers WHERE follower_id = ? AND followed_id = ?"
	result, err := r.DB.Exec(query, followerID, followedID)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no follow relationship found")
	}

	return nil
}

// AcceptFollowRequest updates a follow request status to accepted
func (r *FollowerRepository) AcceptFollowRequest(followerID, followedID string) error {
	query := "UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND followed_id = ? AND status = 'pending'"
	result, err := r.DB.Exec(query, followerID, followedID)
	if err != nil {
		return fmt.Errorf("failed to accept follow request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no pending follow request found")
	}

	return nil
}

// RejectFollowRequest removes a pending follow request
func (r *FollowerRepository) RejectFollowRequest(followerID, followedID string) error {
	query := "DELETE FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'pending'"
	result, err := r.DB.Exec(query, followerID, followedID)
	if err != nil {
		return fmt.Errorf("failed to reject follow request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no pending follow request found")
	}

	return nil
}

// GetPendingFollowRequests retrieves all pending follow requests for a user
func (r *FollowerRepository) GetPendingFollowRequests(userID string) ([]model.FollowerWithUser, error) {
	query := `
		SELECT 
			f.follower_id,
			f.followed_id,
			f.status,
			f.created_at,
			u.id,
			u.fname,
			u.lname,
			u.email,
			COALESCE(u.nickname, '') as nickname,
			COALESCE(u.imgurl, '') as imgurl,
			u.profileVisibility
		FROM followers f
		JOIN users u ON f.follower_id = u.id
		WHERE f.followed_id = ? AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending requests: %w", err)
	}
	defer rows.Close()

	var requests []model.FollowerWithUser
	for rows.Next() {
		var request model.FollowerWithUser
		err := rows.Scan(
			&request.FollowerID,
			&request.FollowedID,
			&request.Status,
			&request.CreatedAt,
			&request.UserID,
			&request.FirstName,
			&request.LastName,
			&request.Email,
			&request.Nickname,
			&request.ImgURL,
			&request.ProfileVisibility,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pending request: %w", err)
		}
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pending requests: %w", err)
	}

	return requests, nil
}
