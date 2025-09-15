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

// GetGroupMembers retrieves all active members of a group.
func (r *GroupRepository) GetGroupMembers(groupID uint) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT user_id FROM group_members
		WHERE group_id = ? AND status = 'active'
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberIDs []string
	for rows.Next() {
		var memberID string
		if err := rows.Scan(&memberID); err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, memberID)
	}

	return memberIDs, rows.Err()
}

// CreateEvent inserts a new event into the database.
func (r *GroupRepository) CreateEvent(event *model.Event) (int, error) {
	query := `
		INSERT INTO events (group_id, creator_id, title, description, event_time)
		VALUES (?, ?, ?, ?, ?)
	`
	res, err := r.DB.Exec(query, event.GroupID, event.CreatorID, event.Title, event.Description, event.EventTime)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// CreateGroupInvitation inserts a new group invitation into the database.
func (r *GroupRepository) CreateGroupInvitation(groupID uint, inviterID, targetUserID string) error {
	_, err := r.DB.Exec(`
		INSERT INTO group_invites (group_id, inviter_user_id, invited_user_id, status)
		VALUES (?, ?, ?, 'pending')
	`, groupID, inviterID, targetUserID)
	return err
}

// AcceptGroupInvitation marks a group invitation as accepted.
func (r *GroupRepository) AcceptGroupInvitation(invitationID int) error {
	_, err := r.DB.Exec(`
		UPDATE group_invites SET status = 'accepted' WHERE id = ?
	`, invitationID)
	return err
}

// DeclineGroupInvitation marks a group invitation as declined.
func (r *GroupRepository) DeclineGroupInvitation(invitationID int) error {
	_, err := r.DB.Exec(`
		UPDATE group_invites SET status = 'declined' WHERE id = ?
	`, invitationID)
	return err
}

// GetGroupInvitation retrieves a group invitation by its ID.
func (r *GroupRepository) GetGroupInvitation(invitationID int) (*model.GroupInvite, error) {
	var invite model.GroupInvite
	err := r.DB.QueryRow(`
		SELECT id, group_id, inviter_user_id, invited_user_id, status
		FROM group_invites
		WHERE id = ?
	`, invitationID).Scan(&invite.ID, &invite.GroupID, &invite.InviterUserID, &invite.InvitedUserID, &invite.Status)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}
