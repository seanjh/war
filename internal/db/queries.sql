-- name: GetSession :one
SELECT id, created FROM sessions
WHERE id = ? LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (id) VALUES (?) RETURNING id, created;

-- name: GetGameSessions :many
SELECT
game_id, COALESCE(session_id, ''), role, hand, battle
FROM game_sessions
WHERE game_id = ?
ORDER BY role;

-- name: GetGameHands :many
SELECT
game_id, card
FROM hands
WHERE game_id = ?;

-- name: CreateHostGameSession :exec
INSERT INTO game_sessions (game_id, session_id, role, hand) VALUES (?, ?, 1, ?), (?, NULL, 2, ?);

-- name: GuestJoinGameSession :exec
UPDATE game_sessions
SET session_id = ?
WHERE game_id = ? AND session_id IS NULL;

-- name: CreateNewHand :exec
INSERT INTO hands (game_session_id, card_slug) VALUES (?, ?);

-- name: GetGame :one
SELECT id, code FROM games
WHERE id = ? LIMIT 1;

-- name: CreateGame :one
INSERT INTO games (id) VALUES (NULL) RETURNING id, code;
