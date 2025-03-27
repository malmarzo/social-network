CREATE TABLE event_options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER,
    option_text TEXT NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id)
);