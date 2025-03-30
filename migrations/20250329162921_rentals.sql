-- +goose Up
-- +goose StatementBegin
CREATE TABLE rentals
(
    parking_lot_id int    NOT NULL,
    start_at       TIMESTAMP,
    end_at         TIMESTAMP,
    cost_per_hour  integer NOT NULL,
    cost_per_day  integer NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rentals;
-- +goose StatementEnd
