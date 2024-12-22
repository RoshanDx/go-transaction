-- name: GetUser :one
SELECT id, username, activated, created_at
FROM users
WHERE id = $1;