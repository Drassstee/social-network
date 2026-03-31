CREATE TABLE IF NOT EXISTS sessions (
    uuid TEXT NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);