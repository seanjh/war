CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL,
    created TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
) STRICT;

CREATE TABLE game_sessions (
    game_id INTEGER NOT NULL,
    session_id INTEGER NOT NULL,
    role TEXT,
    created TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    PRIMARY KEY (game_id, session_id)
) STRICT;
