-- Create a new table without the 'status' column
CREATE TABLE temp_group_chats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id Text NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Copy existing data to the temporary table (ignoring the 'status' column)
INSERT INTO temp_group_chats (id, group_id, user_id, message, created_at)
SELECT id, group_id, user_id, message, created_at FROM group_chats;

-- Drop the modified 'group_chats' table
DROP TABLE group_chats;

-- Rename the temporary table back to 'group_chats'
ALTER TABLE temp_group_chats RENAME TO group_chats;
