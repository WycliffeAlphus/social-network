
CREATE TABLE users (
    id TEXT PRIMARY KEY, -- UUID v4
    email VARCHAR(254) NOT NULL UNIQUE,    
    fname VARCHAR(30) NOT NULL,
    lname VARCHAR(30) NOT NULL,
    dob DATE NOT NULL,
    imgurl VARCHAR(255),
    nickname VARCHAR(30) UNIQUE,
    about TEXT,
    password VARCHAR(255) NOT NULL,
    profileVisibility TEXT NOT NULL DEFAULT 'public' CHECK(profileVisibility IN ('public', 'private')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY, -- UUID v4, the session token
    user_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);