-- +goose Up
CREATE TABLE chirps(
    id UUID not null PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    body varchar(255) not null,
    user_id UUID REFERENCES users ON DELETE CASCADE not null
);

-- +goose Down
DROP TABLE chirps;