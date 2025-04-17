-- +goose Up

-- Create a temporary table with the new column option_id
CREATE TABLE event_options_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    option_text TEXT NOT NULL,
    option_id INTEGER ,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Copy existing data and set option_id as NULL
INSERT INTO event_options_temp (id, event_id, option_text)
SELECT id, event_id, option_text
FROM event_options;

-- Drop the old table
DROP TABLE event_options;

-- Rename the temp table to event_options
ALTER TABLE event_options_temp RENAME TO event_options;
