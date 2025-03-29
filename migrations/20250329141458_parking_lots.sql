-- +goose Up
-- +goose StatementBegin
CREATE TABLE parking_lots
(
    id             TEXT PRIMARY KEY,
    parking_type   TEXT NOT NULL,
    parking_status TEXT NOT NULL,
    owner_id       UUID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE parking_lots;
-- +goose StatementEnd
