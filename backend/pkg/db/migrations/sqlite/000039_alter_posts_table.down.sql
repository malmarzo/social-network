-- Create the original table structure
CREATE TABLE posts_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    privacy TEXT CHECK(privacy IN ('public', 'almost_private', 'private')) DEFAULT 'public',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Copy data back (excluding new columns)
INSERT INTO posts_old (id, user_id, content, image, privacy, created_at)
SELECT id, user_id, content, image, privacy, created_at
FROM posts;

-- Drop the new table
DROP TABLE posts;

-- Rename back to original name
ALTER TABLE posts_old RENAME TO posts;