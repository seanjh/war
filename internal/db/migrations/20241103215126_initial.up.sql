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
    cards_hidden TEXT NOT NULL DEFAULT '',
    created TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    UNIQUE (game_id, session_id)
) STRICT;

CREATE INDEX game_sessions_game_index
ON game_sessions (game_id)
WHERE role IN (1,2);

CREATE TABLE hands (
    id INTEGER PRIMARY KEY,
    game_id INTEGER NOT NULL,
    game_session_id NOT NULL,
    card TEXT NOT NULL,
    position INTEGER NOT NULL CHECK (position BETWEEN 0 AND 52),
    -- cards where "is_hidden = 1" are those won in previous rounds
    is_hidden INTEGER NOT NULL CHECK (is_hidden IN (0,1)) DEFAULT 0,
    created TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_session_id) REFERENCES game_sessions(id),
    FOREIGN KEY (game_id) REFERENCES games(id),
    UNIQUE (game_id, card)
);

CREATE TABLE battles (
    id INTEGER PRIMARY KEY,
    game_session_id NOT NULL,
    round INTEGER NOT NULL,
    card TEXT NOT NULL,
    is_hidden INTEGER NOT NULL CHECK (is_hidden IN (0,1)) DEFAULT 0,
    created TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_session_id) REFERENCES game_sessions(id),
)

-- war
-- 2 card pairs battling
-- N pairs of 2
-- each pair after the first includes 3 additional hidden cards

CREATE INDEX hands_game_index
ON hands (game_id);

CREATE TRIGGER enforce_hands_game_id_session_match
BEFORE INSERT ON hands
BEGIN
    SELECT RAISE(ABORT, 'Mismatched game_id')
    WHERE NEW.game_id != (
        SELECT game_id
        FROM game_sessions
        WHERE id = NEW.game_session_id
    );
END;
