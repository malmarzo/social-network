-- Create table without allowedUsers column
CREATE TABLE posts_old (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    privacy TEXT CHECK(privacy IN ('public', 'almost_private', 'private')) DEFAULT 'public',
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    num_likes INTEGER DEFAULT 0,
    num_dislikes INTEGER DEFAULT 0,
    num_comments INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Copy data back excluding allowedUsers
INSERT INTO posts_old (
    id, user_id, title, content, image,
    privacy, created_at, num_likes,
    num_dislikes, num_comments
)
SELECT 
    id, user_id, title, content, image,
    privacy, created_at, num_likes,
    num_dislikes, num_comments
FROM posts;

-- Drop the current table
DROP TABLE posts;

-- Rename back to original name
ALTER TABLE posts_old RENAME TO posts;