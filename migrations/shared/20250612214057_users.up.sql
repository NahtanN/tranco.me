CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY NOT NULL, 
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    password TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL
);
