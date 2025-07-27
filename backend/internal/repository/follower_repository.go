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
