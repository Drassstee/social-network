-- SQLite doesn't support DROP COLUMN easily before 3.35.0, 
-- but for simplicity in this dev environment we can just leave them or reconstruct if needed.
-- However, standard practice is to provide a down migration.
-- Since we are using golang-migrate, it expects a .down.sql.

-- Note: In SQLite, you usually have to rename table, create new one, copy data.
-- For this exercise, I'll just provide the structure but be aware of SQLite limitations.

PRAGMA foreign_keys=off;

CREATE TABLE messages_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    body TEXT NOT NULL,
    image_url TEXT,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
);
INSERT INTO messages_new SELECT id, sender_id, receiver_id, body, image_url, created_at FROM messages;
DROP TABLE messages;
ALTER TABLE messages_new RENAME TO messages;

CREATE TABLE group_members_new (
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT NOT NULL DEFAULT 'member',
    joined_at DATETIME NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
INSERT INTO group_members_new SELECT group_id, user_id, role, joined_at FROM group_members;
DROP TABLE group_members;
ALTER TABLE group_members_new RENAME TO group_members;

PRAGMA foreign_keys=on;
