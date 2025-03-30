-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS parking_lots
(
    id             INTEGER PRIMARY KEY,
    parking_type   TEXT NOT NULL,
    vehicle_id     TEXT,
    owner_id       UUID
);

-- Insert parking lots with specified distribution
INSERT INTO parking_lots (id, parking_type, vehicle_id, owner_id)
WITH
-- First 100 lots: 80% Permanent (80), 20% Rent (20) in random order
permanent_rent AS (
    SELECT
        generate_series(1, 100) AS id,
        CASE WHEN random() < 0.8
                 THEN 'PERMANENT_PARKING_TYPE'
             ELSE 'RENT_PARKING_TYPE' END AS parking_type
),
-- Next 10 lots: Special
special AS (
    SELECT
        generate_series(101, 110) AS id,
        'SPECIAL_PARKING_TYPE' AS parking_type
),
-- Last 10 lots: Inclusive
inclusive AS (
    SELECT
        generate_series(111, 120) AS id,
        'INCLUSIVE_PARKING_TYPE' AS parking_type
),
-- Combine all lots
all_lots AS (
    SELECT * FROM permanent_rent
    UNION ALL
    SELECT * FROM special
    UNION ALL
    SELECT * FROM inclusive
)
SELECT
    id,
    parking_type,
    CASE WHEN random() < 0.3 THEN gen_random_uuid()::TEXT END AS vehicle_id,
    CASE WHEN random() < 0.3 THEN gen_random_uuid() END AS owner_id
FROM all_lots
ORDER BY id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS parking_lots;
-- +goose StatementEnd