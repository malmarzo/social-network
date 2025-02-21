CREATE TABLE notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type TEXT CHECK(type IN ('follow_request', 'group_invite', 'group_join_request', 'event_created')) NOT NULL,
    reference_id INTEGER NOT NULL, -- Could refer to a user, group, or event
    seen BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
