-- Create new table with allowedUsers column
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
    allowedUsers TEXT, -- Stores comma-separated user IDs
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Copy existing data
INSERT INTO posts_new (
    id, user_id, title, content, image, 
    privacy, created_at, num_likes, 
    num_dislikes, num_comments
)
SELECT 
    id, user_id, title, content, image, 
    privacy, created_at, num_likes, 
    num_dislikes, num_comments
FROM posts;

-- Drop the old table
DROP TABLE posts;

-- Rename the new table
ALTER TABLE posts_new RENAME TO posts;