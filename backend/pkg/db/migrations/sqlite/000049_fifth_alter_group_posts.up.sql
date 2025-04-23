CREATE TABLE group_posts_temp (
    id TEXT PRIMARY KEY,
    group_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    num_likes INTEGER DEFAULT 0,
    num_dislikes INTEGER DEFAULT 0,
    num_comments INTEGER DEFAULT 0,
    created_at TEXT NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO group_posts_temp (
    id, group_id, user_id, user_name, title, content, image, num_likes, num_dislikes, num_comments, created_at
)
SELECT
    id, group_id, user_id, user_name, title, content, image, num_likes, num_dislikes, num_comments, created_at
FROM group_posts;

DROP TABLE group_posts;

ALTER TABLE group_posts_temp RENAME TO group_posts;
