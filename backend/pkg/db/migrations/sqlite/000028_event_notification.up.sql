
CREATE TABLE event_notification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  
    user_id TEXT NOT NULL,
    event_id INTEGER NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'delivered')),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
