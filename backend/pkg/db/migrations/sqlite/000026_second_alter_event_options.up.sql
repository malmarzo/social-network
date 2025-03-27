-- Create a temporary table without the user_id column
CREATE TABLE event_options_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    option_text TEXT NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Copy existing data into the temporary table, excluding the user_id column
INSERT INTO event_options_temp (id, event_id, option_text)
SELECT id, event_id, option_text
FROM event_options;

-- Drop the old table with the user_id column
DROP TABLE event_options;

-- Rename the temporary table to event_options
ALTER TABLE event_options_temp RENAME TO event_options;
