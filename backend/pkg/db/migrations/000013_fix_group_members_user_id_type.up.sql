-- Fix user_id column type in group_members table
-- SQLite doesn't support ALTER COLUMN, so we need to recreate the table

-- Create new table with correct schema
CREATE TABLE group_members_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id VARCHAR(40) NOT NULL,
    role TEXT DEFAULT 'member' NOT NULL,
    status TEXT DEFAULT 'active' NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(group_id, user_id)
);

-- Copy data from old table (if any exists)
INSERT INTO group_members_new (id, group_id, user_id, role, status, joined_at, created_at, updated_at, deleted_at)
SELECT id, group_id, CAST(user_id AS TEXT), role, status, joined_at, created_at, updated_at, deleted_at
FROM group_members
WHERE EXISTS (SELECT 1 FROM group_members LIMIT 1);

-- Drop old table
DROP TABLE group_members;

-- Rename new table
ALTER TABLE group_members_new RENAME TO group_members;
