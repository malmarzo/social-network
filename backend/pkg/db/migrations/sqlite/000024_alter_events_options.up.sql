
CREATE TABLE event_options_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    option_text TEXT NOT NULL,
    user_id Text NOT NULL, -- Added new user_id column
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);


INSERT INTO event_options_temp (id, event_id, option_text, user_id)
SELECT id, event_id, option_text, 0 FROM event_options;


DROP TABLE event_options;


ALTER TABLE event_options_temp RENAME TO event_options;
