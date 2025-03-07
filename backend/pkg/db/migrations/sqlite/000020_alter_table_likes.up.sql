CREATE TABLE likes_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    type TEXT CHECK(type IN ('like', 'dislike', 'none')) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    UNIQUE(user_id, post_id) -- Prevent duplicate likes/dislikes
);

INSERT INTO likes_new (id, user_id, post_id, type)
SELECT id, user_id, post_id, type
FROM likes;

DROP TABLE likes;

ALTER TABLE likes_new RENAME TO likes;
