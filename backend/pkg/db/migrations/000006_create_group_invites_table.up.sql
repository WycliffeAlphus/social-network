CREATE TABLE group_invites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    inviter_user_id INTEGER NOT NULL,
    invited_user_id INTEGER NOT NULL,
    status TEXT DEFAULT 'pending' NOT NULL, -- 'pending', 'accepted', 'rejected', 'expired'
    invited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL, -- For soft deletes
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (inviter_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (invited_user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(group_id, invited_user_id, status) -- Prevent duplicate pending invites for the same group and user
);