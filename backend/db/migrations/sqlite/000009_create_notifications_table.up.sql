CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    actor_id INTEGER NOT NULL,
    target_type TEXT NOT NULL, -- 'post', 'comment', 'group_invite', 'group_request', etc.
    target_id INTEGER NOT NULL,
    notification_type TEXT NOT NULL, -- 'like', 'dislike', 'comment', 'invite', 'request', 'accept', 'decline'
    is_read BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (actor_id) REFERENCES users(id)
);
