-- name: CreateTodo :one
INSERT INTO todo (
    account_id,
    title,
    time,
    date,
    complete
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTodo :one
SELECT * FROM todo
WHERE id = $1 LIMIT 1;

-- name: ListTodo :many
SELECT * FROM todo
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTodo :one
UPDATE todo 
SET title = $2, time = $3, date = $4, complete = $5
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todo
WHERE id = $1;