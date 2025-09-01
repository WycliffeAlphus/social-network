package repository

import (
	"database/sql"

	"backend/internal/model"
)

type GroupRepository struct {
	DB *sql.DB 
}

// NewGroupRepository creates and returns a new instance of GroupRepository.
func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{DB: db}
}

// FindAll retrieves all groups from the database.
func (r *GroupRepository) FindAll() ([]model.Group, error) {
	rows, err := r.DB.Query("SELECT id, title, description, creator_id, privacy_setting, created_at FROM groups ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var group model.Group
		if err := rows.Scan(&group.ID, &group.Title, &group.Description, &group.CreatorID, &group.PrivacySetting, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// InsertGroup inserts a new group into the database using a transaction.
func (r *GroupRepository) InsertGroup(tx *sql.Tx, group *model.Group) (uint, error) {
	stmt, err := tx.Prepare(`INSERT INTO groups (title, description, creator_id, privacy_setting) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(group.Title, group.Description, group.CreatorID, group.PrivacySetting)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// InsertGroupMember inserts a new group member record using a transaction.
func (r *GroupRepository) InsertGroupMember(tx *sql.Tx, member *model.GroupMember) error {
	stmt, err := tx.Prepare(`INSERT INTO group_members (group_id, user_id, role, status) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(member.GroupID, member.UserID, member.Role, member.Status)
	return err
}

// FindGroupByID retrieves a group by its ID.
func (r *GroupRepository) FindGroupByID(groupID uint) (*model.Group, error) {
	var group model.Group
	err := r.DB.QueryRow(`
		SELECT id, title, description, creator_id, privacy_setting, created_at, updated_at
		FROM groups
		WHERE id = ? AND deleted_at IS NULL
	`, groupID).Scan(&group.ID, &group.Title, &group.Description, &group.CreatorID, &group.PrivacySetting, &group.CreatedAt, &group.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Group not found
		}
		return nil, err
	}
	return &group, nil
}

// CheckUserMembership checks if a user is already a member of a group.
func (r *GroupRepository) CheckUserMembership(groupID uint, userID string) (bool, string, error) {
	var status string
	err := r.DB.QueryRow(`
		SELECT status
		FROM group_members
		WHERE group_id = ? AND user_id = ? AND deleted_at IS NULL
	`, groupID, userID).Scan(&status)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil // User is not a member
		}
		return false, "", err
	}
	return true, status, nil
}

// CreateJoinRequest creates a pending membership request for a user to join a group.
func (r *GroupRepository) CreateJoinRequest(groupID uint, userID string) error {
	_, err := r.DB.Exec(`
		INSERT INTO group_members (group_id, user_id, role, status)
		VALUES (?, ?, 'member', 'pending')
	`, groupID, userID)
	return err
}

// AcceptJoinRequest updates a pending join request to active status.
func (r *GroupRepository) AcceptJoinRequest(groupID uint, userID string) error {
	result, err := r.DB.Exec(`
		UPDATE group_members
		SET status = 'active', updated_at = CURRENT_TIMESTAMP
		WHERE group_id = ? AND user_id = ? AND status = 'pending' AND deleted_at IS NULL
	`, groupID, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // No pending request found
	}

	return nil
}

// IsGroupCreator checks if a user is the creator of a group.
func (r *GroupRepository) IsGroupCreator(groupID uint, userID string) (bool, error) {
	var creatorID string
	err := r.DB.QueryRow(`
		SELECT creator_id
		FROM groups
		WHERE id = ? AND deleted_at IS NULL
	`, groupID).Scan(&creatorID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Group not found
		}
		return false, err
	}

	return creatorID == userID, nil
}

// GetPendingJoinRequests retrieves all pending join requests for a group.
func (r *GroupRepository) GetPendingJoinRequests(groupID uint) ([]model.GroupJoinRequest, error) {
	rows, err := r.DB.Query(`
		SELECT gm.user_id, u.fname, u.lname, u.imgurl, gm.created_at
		FROM group_members gm
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = ? AND gm.status = 'pending' AND gm.deleted_at IS NULL
		ORDER BY gm.created_at ASC
	`, groupID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []model.GroupJoinRequest
	for rows.Next() {
		var req model.GroupJoinRequest
		var imgURL sql.NullString

		err := rows.Scan(&req.UserID, &req.FirstName, &req.LastName, &imgURL, &req.RequestedAt)
		if err != nil {
			return nil, err
		}

		req.GroupID = groupID
		req.UserName = req.FirstName + " " + req.LastName
		if imgURL.Valid {
			req.UserImageURL = imgURL.String
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// RejectJoinRequest removes a pending join request.
func (r *GroupRepository) RejectJoinRequest(groupID uint, userID string) error {
	result, err := r.DB.Exec(`
		DELETE FROM group_members
		WHERE group_id = ? AND user_id = ? AND status = 'pending' AND deleted_at IS NULL
	`, groupID, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // No pending request found
	}

	return nil
}
