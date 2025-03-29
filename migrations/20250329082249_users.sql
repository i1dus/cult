-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone      TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    surname    TEXT NOT NULL,
    patronymic TEXT,
    address    TEXT,
    user_type  TEXT NOT NULL,
    pass_hash  TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
