
CREATE TABLE events_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    creator_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    event_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (creator_id) REFERENCES users(id)
);


INSERT INTO events_temp (id, group_id, creator_id, title, description, event_date, created_at)
SELECT id, group_id, CAST(creator_id AS INTEGER), title, description, event_date, created_at
FROM events;


DROP TABLE events;
ALTER TABLE events_temp RENAME TO events;
