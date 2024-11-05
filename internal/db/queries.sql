-- name: GetSession :one
SELECT id, created FROM sessions
WHERE id = ? LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (id) VALUES (?) RETURNING id, created;

-- name: GetGameSession :many
SELECT s.game_id, s.session_id, g.code
FROM game_sessions s
INNER JOIN games g ON g.id = s.game_id
WHERE s.game_id = ?
ORDER BY s.session_id;

-- name: CreateGameSession :one
INSERT INTO game_sessions (game_id, session_id, role) VALUES (?, ?, ?) RETURNING game_id, session_id, role;

-- name: GetGame :one
SELECT id, code FROM games
WHERE id = ? LIMIT 1;

-- name: CreateGame :one
INSERT INTO games (id, code) VALUES (NULL, NULL) RETURNING id, code;

-- name: CreateDeck :exec
INSERT INTO decks (game_id, session_id, cards) VALUES (?, ?, ?);
