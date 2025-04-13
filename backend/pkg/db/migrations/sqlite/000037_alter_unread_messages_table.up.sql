
CREATE TABLE unread_messages_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    count INTEGER NOT NULL DEFAULT 0,
    group_message_id INTEGER NOT NULL, -- New column added here
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_message_id) REFERENCES group_chats(id) ON DELETE CASCADE -- Foreign key for group_message_id
);


INSERT INTO unread_messages_temp (
    id, group_id, user_id, count, group_message_id
)
SELECT
    id, group_id, user_id, count, 0 AS group_message_id 
FROM unread_messages;


DROP TABLE unread_messages;


ALTER TABLE unread_messages_temp RENAME TO unread_messages;
