CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    fname TEXT NOT NULL,
    lname TEXT NOT NULL,
    dob TEXT NOT NULL,
    imgurl TEXT,
    nickname TEXT,
    about TEXT,
    password TEXT NOT NULL,
    profileVisibility TEXT DEFAULT 'public',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
