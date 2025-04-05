CREATE TABLE comments_old (
    id INT PRIMARY KEY AUTOINCREMENT,
    user_id INT NOT NULL,
    post_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
)
    
INSERT INTO comments_old (user_id, post_id, content, created_at)
SELECT user_id, post_id, content, created_at
FROM comments;

DROP TABLE comments;

ALTER TABLE comments_old RENAME TO comments;