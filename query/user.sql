-- name: GetAllUser :many
SELECT * FROM users;

-- name: GetUser :one
SELECT id, username, activated
FROM users
WHERE id = $1;

-- name: InsertUser :one
INSERT INTO users (username, firstname, activated)
VALUES ($1, $2, $3)
RETURNING *;
