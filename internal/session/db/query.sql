-- name: GetSession :one
SELECT id, created FROM sessions
WHERE id = ? LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (id) VALUES (NULL) RETURNING id, created;
