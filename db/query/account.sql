-- name: CreateAccount :one
INSERT INTO account (
    first_name,
    last_name,
    email,
    user_name,
    password
) VALUES (
    $1, $2, $3, $4, $5
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
SET 
    first_name = COALESCE(sqlc.narg(first_name),first_name),
    last_name = COALESCE(sqlc.narg(last_name),last_name),
    email = COALESCE(sqlc.narg(email),email),
    password = COALESCE(sqlc.narg(password),password)
WHERE 
    user_name = sqlc.arg(user_name)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;