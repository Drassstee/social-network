-- Add is_read to private messages
ALTER TABLE messages ADD COLUMN is_read BOOLEAN NOT NULL DEFAULT 0;

-- Add last_read_msg_id to group_members to track per-user group chat progress
ALTER TABLE group_members ADD COLUMN last_read_msg_id INTEGER NOT NULL DEFAULT 0;
