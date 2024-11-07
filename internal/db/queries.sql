-- name: GetSession :one
SELECT id, created FROM sessions
WHERE id = ? LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (id) VALUES (?) RETURNING id, created;

-- name: GetGameSessions :many
SELECT game_id, COALESCE(session_id, ''), role, deck
FROM game_sessions
WHERE game_id = ?
ORDER BY role;

-- name: CreateHostGameSession :exec
INSERT INTO game_sessions (game_id, session_id, role, deck) VALUES (?, ?, 1, ?), (?, NULL, 2, ?);

-- name: GetGame :one
SELECT id, code FROM games
WHERE id = ? LIMIT 1;

-- name: CreateGame :one
INSERT INTO games (id) VALUES (NULL) RETURNING id, code;
