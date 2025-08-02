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
