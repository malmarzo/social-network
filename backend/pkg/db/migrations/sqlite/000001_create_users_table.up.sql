CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    avatar TEXT, -- Image path
    nickname TEXT UNIQUE,
    about_me TEXT,
    is_private BOOLEAN DEFAULT 0, -- 0 = Public, 1 = Private
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
