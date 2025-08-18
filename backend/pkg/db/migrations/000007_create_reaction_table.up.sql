CREATE TABLE reactions (
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    type TEXT CHECK(type IN ('like', 'dislike')) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);