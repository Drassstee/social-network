CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    password TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    avatar TEXT,
    nickname TEXT,
    about_me TEXT,
    profile_type TEXT DEFAULT 'public'
);
