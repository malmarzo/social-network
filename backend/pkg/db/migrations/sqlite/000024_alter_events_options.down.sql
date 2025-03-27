
CREATE TABLE event_options_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    option_text TEXT NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id)
);



INSERT INTO event_options_temp (id, event_id, option_text)
SELECT id, event_id, option_text FROM event_options;


DROP TABLE event_options;

ALTER TABLE event_options_temp RENAME TO event_options;
