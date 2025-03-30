-- +goose Up
-- +goose StatementBegin
CREATE TABLE bookings
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rental_id     UUID     NOT NULL,
    user_id       UUID    NOT NULL,
    vehicle       TEXT    NOT NULL,
    is_short_term BOOLEAN NOT NULL,
    is_present    BOOLEAN NOT NULL DEFAULT TRUE,
    start_at      timestamp,
    end_at        timestamp,
    UNIQUE (rental_id, start_at, end_at)
);

INSERT INTO bookings (rental_id, user_id, vehicle, is_short_term, is_present, start_at, end_at)
VALUES
-- Current active bookings (is_present = TRUE)
('1a2b3c4d-5e6f-7890-1234-567890abcdef', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'А123БВ78', TRUE, TRUE, '2024-06-15 08:00:00', '2025-06-15 12:30:00'),
('2b3c4d5e-6f7a-8901-2345-67890abcdef1', '550e8400-e29b-41d4-a716-446655440000', 'Е456КХ78', FALSE, TRUE, '2023-06-15 09:15:00', '2025-06-15 18:00:00'),
('3c4d5e6f-7a8b-9012-3456-7890abcdef12', '6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'О789РТ78', TRUE, TRUE, '2023-06-15 10:30:00', '2023-06-15 14:45:00'),

-- Future bookings (is_present = FALSE)
('4d5e6f7a-8b9c-0123-4567-89abcdef1234', '6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'У321АВ78', TRUE, FALSE, '2025-06-16 09:00:00', '2026-06-16 11:30:00'),
('5e6f7a8b-9c0d-1234-5678-9abcdef12345', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Т654СН78', FALSE, FALSE, '2023-06-17 08:00:00', '2025-06-17 17:00:00'),

-- Completed bookings (in the past)
('6f7a8b9c-0d1e-2345-6789-abcdef123456', '550e8400-e29b-41d4-a716-446655440000', 'М987КХ78', TRUE, FALSE, '2023-06-14 10:00:00', '2025-06-14 12:00:00'),
('7a8b9c0d-1e2f-3456-789a-bcdef1234567', '6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'Х159УК78', FALSE, FALSE, '2023-06-13 09:00:00', '2025-06-13 18:00:00');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bookings;
-- +goose StatementEnd
