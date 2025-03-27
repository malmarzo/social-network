
CREATE TABLE event_participation_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT CHECK(status IN ('going', 'not_going')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(event_id, user_id)
);


INSERT INTO event_participation_temp (id, event_id, user_id)
SELECT id, event_id, user_id
FROM event_participation;


DROP TABLE event_participation;


ALTER TABLE event_participation_temp RENAME TO event_participation;
