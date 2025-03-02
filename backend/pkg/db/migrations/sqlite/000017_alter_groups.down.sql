CREATE TABLE groups_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    creator_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

INSERT INTO groups_old (id, creator_id, title, description, created_at)
SELECT id, CAST(creator_id AS INTEGER), title, description, created_at FROM groups;

DROP TABLE groups;
ALTER TABLE groups_old RENAME TO groups;

