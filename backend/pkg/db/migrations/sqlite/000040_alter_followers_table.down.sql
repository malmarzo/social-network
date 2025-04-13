CREATE TABLE followers_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    following_id INTEGER NOT NULL,
    status TEXT CHECK( status IN ('pending', 'accepted') ) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (following_id) REFERENCES users(id),
    UNIQUE(follower_id, following_id)
);

INSERT INTO followers_old (id, follower_id, following_id, status)
SELECT id, follower_id, following_id, status
FROM followers;

DROP TABLE followers;   

ALTER TABLE followers_old RENAME TO followers;