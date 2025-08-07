package repository

import (
	"backend/internal/model"
	"database/sql"
)

func CreateReaction(reaction *model.Reaction, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO reactions (post_id, user_id, type) VALUES (?, ?, ?)",
		reaction.PostID, reaction.UserID, reaction.Type)
	return err
}

func CheckIfReactionExist(reaction *model.Reaction, db *sql.DB) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM reactions WHERE post_id = ? AND user_id = ? AND type = ?)",
		reaction.PostID, reaction.UserID, reaction.Type).Scan(&exists)
	return exists, err
}
func GetReactionsCount(postID string, db *sql.DB) (likes int, dislikes int, err error) {
	err = db.QueryRow("SELECT COUNT(*) FROM reactions WHERE post_id = ? AND type = 'like'", postID).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}
	err = db.QueryRow("SELECT COUNT(*) FROM reactions WHERE post_id = ? AND type = 'dislike'", postID).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}

func DeleteReaction(reaction *model.Reaction, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM reactions WHERE post_id = ? AND user_id = ? AND type = ?",
		reaction.PostID, reaction.UserID, reaction.Type)
	return err
}

func UpdateReaction(reaction *model.Reaction, db *sql.DB) error {
	_, err := db.Exec("UPDATE reactions SET type = ? WHERE post_id = ? AND user_id = ?",
		reaction.Type, reaction.PostID, reaction.UserID)
	return err
}

func CheckIfUserAlreadyReacted(reaction *model.Reaction, db *sql.DB) (bool, error)  {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM reactions WHERE post_id = ? AND user_id = ? )",
		reaction.PostID, reaction.UserID).Scan(&exists)
	return exists, err
}