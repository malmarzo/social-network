CREATE TABLE followers_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id TEXT NOT NULL,
    following_id TEXT NOT NULL,
    status TEXT CHECK( status IN ('pending', 'accepted') ) DEFAULT 'pending',
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (following_id) REFERENCES users(id),
    UNIQUE(follower_id, following_id)
);

INSERT INTO followers_new (id, follower_id, following_id, status)
SELECT id, follower_id, following_id, status
FROM followers;

DROP TABLE followers;

ALTER TABLE followers_new RENAME TO followers;