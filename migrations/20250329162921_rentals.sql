-- +goose Up
-- +goose StatementBegin
CREATE TABLE rentals
(
    parking_lot_id TEXT    NOT NULL,
    start_at       TIMESTAMP,
    end_at         TIMESTAMP,
    cost_per_hour  integer NOT NULL,
    created_at     TIMESTAMP DEFAULT now(),
    updated_at     TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rentals;
-- +goose StatementEnd
