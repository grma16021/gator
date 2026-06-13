-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
values (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1; 


-- name: DeleteUsers :exec
DELETE FROM users;


-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, name, url, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;
