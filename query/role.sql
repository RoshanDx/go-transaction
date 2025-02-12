-- name: GetRole :one
SELECT * FROM roles
WHERE name = $1;

-- name: AssignUserRole :exec
INSERT INTO user_role(user_id, role_id)
VALUES ($1, $2);