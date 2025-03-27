CREATE TABLE event_participation_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER,
    user_id INTEGER, -- Revert user_id back to INTEGER
    option_id INTEGER,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (option_id) REFERENCES event_options(id),
    UNIQUE(event_id, user_id)
);

INSERT INTO event_participation_temp (id, event_id, user_id, option_id)
SELECT id, event_id, CAST(user_id AS INTEGER), option_id
FROM event_participation;

DROP TABLE event_participation;

ALTER TABLE event_participation_temp RENAME TO event_participation;
