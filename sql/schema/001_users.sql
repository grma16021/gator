-- +goose Up
Create TABLE users(
    id TEXT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE feeds(
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT,
    url TEXT UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT feeds_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE users, feeds;


