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
    id INTEGER PRIMARY KEY,
    game_id INTEGER NOT NULL,
    -- Games are created by the host (role = 1) with a session_id with a placeholder
    -- guest game_session (session_id = NULL AND role = 2).
    session_id TEXT CHECK (role != 1 OR session_id IS NOT NULL),
    -- Game host is "role = 1" and guest is "role = 2"
    role INTEGER NOT NULL CHECK (role IN (1, 2)),
    cards_in_hand TEXT NOT NULL,
    cards_in_battle TEXT NOT NULL DEFAULT '',
    cards_in_battle_hidden TEXT NOT NULL DEFAULT '',
    cards_won TEXT NOT NULL DEFAULT '',
    created TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    UNIQUE (game_id, session_id)
) STRICT;

CREATE INDEX game_sessions_game_index
ON game_sessions (game_id)
WHERE role IN (1,2);
