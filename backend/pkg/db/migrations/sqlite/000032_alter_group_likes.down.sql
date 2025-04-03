CREATE TABLE group_likes_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL, 
    post_id INTEGER NOT NULL,
    type TEXT CHECK(type IN ('like', 'dislike', 'none')) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id), -- Revert to original reference
    UNIQUE(user_id, post_id) 
);

INSERT INTO group_likes_temp (id, user_id, post_id, type)
SELECT id, user_id, post_id, type FROM group_likes;

DROP TABLE group_likes;

ALTER TABLE group_likes_temp RENAME TO group_likes;
