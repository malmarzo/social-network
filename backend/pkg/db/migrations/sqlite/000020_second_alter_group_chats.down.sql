
CREATE TABLE temp_group_chats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status TEXT CHECK(status IN ('pending', 'delivered')) DEFAULT 'pending',
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);


INSERT INTO temp_group_chats (id, group_id, user_id, message, created_at, status)
SELECT id, group_id, user_id, message, created_at, status FROM group_chats;


DROP TABLE group_chats;


ALTER TABLE temp_group_chats RENAME TO group_chats;
