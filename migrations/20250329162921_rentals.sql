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

INSERT INTO rentals (id, parking_lot_id, start_at, end_at, cost_per_hour, cost_per_day)
VALUES
-- Current rentals (active now)
('1a2b3c4d-5e6f-7890-1234-567890abcdef', 1, NOW() - INTERVAL '1 hour', NOW() + INTERVAL '3 hours', 100, 1500),
('2b3c4d5e-6f7a-8901-2345-67890abcdef1', 2, NOW() - INTERVAL '30 minutes', NOW() + INTERVAL '2 hours', 120, 1800),
('3c4d5e6f-7a8b-9012-3456-7890abcdef12', 3, NOW() - INTERVAL '2 hours', NOW() + INTERVAL '4 hours', 90, 1300),

-- Future bookings
('4d5e6f7a-8b9c-0123-4567-89abcdef1234', 4, NOW() - INTERVAL '1 day', NOW() + INTERVAL '1 day 6 hours', 110, 1600),
('5e6f7a8b-9c0d-1234-5678-9abcdef12345', 5, NOW() - INTERVAL '2 days', NOW() + INTERVAL '2 days 8 hours', 150, 2200),
('6f7a8b9c-0d1e-2345-6789-abcdef123456',6, NOW() - INTERVAL '3 days', NOW() + INTERVAL '3 days 12 hours', 95, 1400),

-- Past rentals (completed)
('7a8b9c0d-1e2f-3456-789a-bcdef1234567', 7, NOW() - INTERVAL '3 days', NOW() + INTERVAL '2 days', 80, 1200);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rentals;
-- +goose StatementEnd
