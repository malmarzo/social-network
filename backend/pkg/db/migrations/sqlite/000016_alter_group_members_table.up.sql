CREATE TABLE new_group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    invited_by TEXT NOT NULL,
    status TEXT CHECK(status IN ('pending', 'accepted', 'declined')) DEFAULT 'pending',
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (invited_by) REFERENCES users(id)
);

INSERT INTO new_group_members (id, group_id, user_id, status)
SELECT id, group_id, CAST(user_id AS TEXT), status FROM group_members;

DROP TABLE group_members;

ALTER TABLE new_group_members RENAME TO group_members;
