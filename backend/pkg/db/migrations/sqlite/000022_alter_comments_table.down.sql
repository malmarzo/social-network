CREATE TABLE old_comments (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    image TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

INSERT INTO new_comments (id, user_id, post_id, content, created_at)
SELECT id, user_id, post_id, content, created_at
FROM comments;

DROP TABLE comments;

ALTER TABLE old_comments RENAME TO comments;