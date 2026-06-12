-- +goose Up
Create TABLE users(
    id int UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(50) NOT NULL
);

-- +goose Down
DROP TABLE users;