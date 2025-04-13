CREATE TABLE likes_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    type TEXT CHECK(type IN ('like', 'dislike')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    UNIQUE(user_id, post_id) -- Prevent duplicate likes/dislikes
);

INSERT INTO likes_old (id, user_id, post_id, type)
SELECT id, user_id, post_id, type
FROM likes; 

DROP TABLE likes;
ALTER TABLE likes_old RENAME TO likes;