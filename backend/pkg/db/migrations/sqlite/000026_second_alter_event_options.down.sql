-- Create a temporary table with the user_id column again
CREATE TABLE event_options_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    option_text TEXT NOT NULL,
    user_id TEXT NOT NULL, -- Re-adding user_id
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Copy existing data back into the temporary table with user_id set to an empty string
INSERT INTO event_options_temp (id, event_id, option_text, user_id)
SELECT id, event_id, option_text, '' FROM event_options;

-- Drop the current table
DROP TABLE event_options;

-- Rename the temporary table to event_options
ALTER TABLE event_options_temp RENAME TO event_options;
