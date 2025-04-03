CREATE TABLE group_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL, 
    post_id INTEGER NOT NULL,
    type TEXT CHECK(type IN ('like', 'dislike', 'none')) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    UNIQUE(user_id, post_id) 
);
