
CREATE TABLE unread_messages_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    count INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


INSERT INTO unread_messages_temp (
    id, group_id, user_id, count
)
SELECT
    id, group_id, user_id, count
FROM unread_messages;

DROP TABLE unread_messages;


ALTER TABLE unread_messages_temp RENAME TO unread_messages;
