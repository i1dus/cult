-- +goose Up
-- +goose StatementBegin
CREATE TABLE rentals
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parking_lot_id INTEGER NOT NULL,
    start_at       TIMESTAMP,
    end_at         TIMESTAMP,
    cost_per_hour  INTEGER NOT NULL,
    cost_per_day   INTEGER NOT NULL,
    UNIQUE (parking_lot_id, start_at, end_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rentals;
-- +goose StatementEnd
