-- +goose Up
-- +goose StatementBegin
CREATE TABLE bookings
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rental_id     int     NOT NULL,
    user_id       UUID    NOT NULL,
    vehicle       TEXT    NOT NULL,
    is_short_term BOOLEAN NOT NULL,
    start_at      timestamp,
    end_at        timestamp,
    UNIQUE (rental_id, start_at, end_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bookings;
-- +goose StatementEnd
