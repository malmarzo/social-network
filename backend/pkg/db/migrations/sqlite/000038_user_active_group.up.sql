CREATE TABLE user_active_group (
    user_id TEXT PRIMARY KEY,
    group_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
