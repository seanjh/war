-- name: GetSession :one
SELECT id, created FROM sessions
WHERE id = ? LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (id) VALUES (NULL) RETURNING id, created;

-- name: GetGameSession :many
SELECT s.game_id, s.session_id, g.code
FROM game_sessions s
INNER JOIN games g ON g.id = s.game_id
WHERE s.game_id = ?
ORDER BY s.session_id;

-- name: CreateGameSession :exec
INSERT INTO game_sessions (game_id, session_id, role) VALUES (?, ?, ?) RETURNING game_id, session_id, role;

-- name: CreateGame :exec
INSERT INTO games (code) VALUES (?) RETURNING id, code;
