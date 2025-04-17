-- +goose Down

-- Create temporary table WITHOUT AUTOINCREMENT on id
CREATE TABLE event_participation_temp (
    id INTEGER PRIMARY KEY,
    event_id INTEGER,
    user_id TEXT,
    option_id INTEGER,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (option_id) REFERENCES event_options(id),
    UNIQUE(event_id, user_id)
);

-- Copy data from current table
INSERT INTO event_participation_temp (id, event_id, user_id, option_id)
SELECT id, event_id, user_id, option_id
FROM event_participation;

-- Drop current table
DROP TABLE event_participation;

-- Rename temporary table to original
ALTER TABLE event_participation_temp RENAME TO event_participation;
