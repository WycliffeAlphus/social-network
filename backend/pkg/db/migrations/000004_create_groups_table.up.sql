CREATE TABLE groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    creator_id INTEGER NOT NULL,
    privacy_setting TEXT DEFAULT 'private' NOT NULL, -- 'public', 'private', 'secret'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL, -- For soft deletes
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);