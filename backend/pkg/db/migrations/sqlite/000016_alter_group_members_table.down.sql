CREATE TABLE group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT CHECK(status IN ('invited', 'member')) DEFAULT 'invited',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(group_id, user_id)
);

INSERT INTO group_members (id, group_id, user_id, status)
SELECT id, group_id, CAST(user_id AS INTEGER), status FROM group_members;

DROP TABLE group_members;

ALTER TABLE group_members RENAME TO group_members;
