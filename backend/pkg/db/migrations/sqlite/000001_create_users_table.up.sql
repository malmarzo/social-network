CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth TEXT NOT NULL,
    avatar TEXT, -- Image path
    nickname TEXT UNIQUE,
    about_me TEXT,
    is_private BOOLEAN DEFAULT 0, -- 0 = Public, 1 = Private
    created_at TEXT DEFAULT CURRENT_TIMESTAMP -- Set default value to current timestamp
);
