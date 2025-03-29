-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS parking_lots
(
    id             INTEGER PRIMARY KEY,
    parking_type   TEXT NOT NULL,
    parking_status TEXT NOT NULL,
    vehicle_id     UUID,
    owner_id       UUID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE parking_lots;
-- +goose StatementEnd
