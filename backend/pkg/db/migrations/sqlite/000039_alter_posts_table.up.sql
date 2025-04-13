-- First, create a new table with the desired structure
CREATE TABLE posts_new (
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

-- Copy existing data (mapping old columns to new ones)
INSERT INTO posts_new (id, user_id, content, image, privacy, created_at)
SELECT id, user_id, content, image, privacy, created_at
FROM posts;

-- Drop the old table
DROP TABLE posts;

-- Rename the new table to posts
ALTER TABLE posts_new RENAME TO posts;