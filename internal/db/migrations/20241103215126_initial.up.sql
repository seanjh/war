CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    created TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
) STRICT;

CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL DEFAULT (hex(randomblob(4))),
    created TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
) STRICT;

CREATE TABLE game_sessions (
    game_id INTEGER NOT NULL,
    session_id TEXT CHECK (role != 1 OR session_id IS NOT NULL),
    role INTEGER NOT NULL CHECK (role IN (1, 2)),
    deck TEXT NOT NULL,
    created TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    UNIQUE (game_id, session_id)
) STRICT;
