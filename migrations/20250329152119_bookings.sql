-- +goose Up
-- +goose StatementBegin
CREATE TABLE bookings
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parking_lot_id int  NOT NULL,
    user_id        UUID NOT NULL,
    vehicle_id     text NOT NULL,
    start_at       timestamp,
    end_at         timestamp,
    UNIQUE (parking_lot_id, start_at, end_at)

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bookings;
-- +goose StatementEnd
