-- +goose Up
CREATE TABLE users(
    id UUID not null PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    email varchar(255) not null UNIQUE
);

-- +goose Down
DROP TABLE users;