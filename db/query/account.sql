-- name: CreateAccount :one
INSERT INTO account (
    first_name,
    last_name,
    user_name,
    password
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM account
WHERE user_name = $1 LIMIT 1;

-- name: ListAccount :many
SELECT * FROM account
WHERE user_name = $1
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE account 
SET first_name = $2, last_name = $3, user_name = $4, password = $5
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;