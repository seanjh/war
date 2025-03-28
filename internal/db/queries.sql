-- name: GetSession :one
SELECT id, created FROM sessions
WHERE id = ? LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (id) VALUES (?) RETURNING id, created;

-- name: GetGameSessions :many
SELECT
game_id, COALESCE(session_id, ''), role, cards_in_hand, cards_in_battle, cards_in_battle_hidden, cards_in_reserve
FROM game_sessions
WHERE game_id = ?
ORDER BY role;

-- name: GetGameSession :one
SELECT
id, game_id, COALESCE(session_id, ''), role, cards_in_hand, cards_in_battle, cards_in_battle_hidden, cards_in_reserve
FROM game_sessions
WHERE game_id = ? AND session_id = ?
LIMIT 1;

-- name: CreateHostGameSession :exec
INSERT INTO game_sessions (game_id, session_id, role, cards_in_hand) VALUES (?, ?, 1, ?), (?, NULL, 2, ?);

-- name: GuestJoinGameSession :exec
UPDATE game_sessions
SET session_id = ?
WHERE game_id = ? AND session_id IS NULL;

-- name: GetGame :one
SELECT id, code FROM games
WHERE id = ? LIMIT 1;

-- name: CreateGame :one
INSERT INTO games (id) VALUES (NULL) RETURNING id, code;

-- name: FlipCard :exec
UPDATE game_sessions
SET cards_in_hand = ?, cards_in_battle = ?, updated = CURRENT_TIMESTAMP
WHERE id = ?;
