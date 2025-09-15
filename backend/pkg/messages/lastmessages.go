package messages

import "database/sql"

// GetLastMessageTimestamp gets the timestamp of the last message between user1 and user2
func GetLastMessageTimestamp(db *sql.DB, userId1, userId2 string) (string, error) {
	var timestamp string
	err := db.QueryRow(`
        SELECT timestamp 
        FROM messages 
        WHERE (from_user = ? AND to_user = ?) OR (from_user = ? AND to_user = ?)
        ORDER BY timestamp DESC
        LIMIT 1`,
		userId1, userId2, userId2, userId1).Scan(&timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return timestamp, nil
}
