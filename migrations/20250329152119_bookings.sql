-- +goose Up
-- +goose StatementBegin
CREATE TABLE bookings
(
    id UUID PRIMARY KEY,
    parking_lot_id TEXT NOT NULL,
    user_id UUID NOT NULL,
    vehicle_id     UUID NOT NULL,
    start_at           timestamp,
    end_at             timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bookings;
-- +goose StatementEnd
