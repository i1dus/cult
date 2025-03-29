-- +goose Up
-- +goose StatementBegin
CREATE TABLE bookings
(
    parking_lot_id TEXT NOT NULL,
    vehicle_id     UUID NOT NULL,
    start_at           timestamp,
    end_at             timestamp,
    PRIMARY KEY (parking_lot_id, vehicle_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bookings;
-- +goose StatementEnd
