
CREATE TABLE group_posts_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);


INSERT INTO group_posts_temp (id, group_id, user_id, content, image, created_at)
SELECT id, group_id, CAST(user_id AS INTEGER), content, image, created_at FROM group_posts;


DROP TABLE group_posts;


CREATE TABLE group_posts AS SELECT * FROM group_posts_temp;
DROP TABLE group_posts_temp;
