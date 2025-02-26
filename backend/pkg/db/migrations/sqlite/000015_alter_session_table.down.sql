CREATE TABLE sessions_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token TEXT NOT NULL,
    expiration TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(session_token)
);

INSERT INTO sessions_new (id, user_id, session_token, expiration,created_at)
SELECT id, user_id, session_token, expiration,created_at FROM sessions;

DROP Table sessions;

ALTER TABLE sessions_new RENAME TO sessions;