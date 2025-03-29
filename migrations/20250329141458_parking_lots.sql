-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS parking_lots
(
    id             TEXT PRIMARY KEY,
    parking_type   TEXT NOT NULL,
    vehicle_id     UUID,
    owner_id       UUID
);

-- Insert parking lots with specified distribution
INSERT INTO parking_lots (id, parking_type, vehicle_id, owner_id)
WITH
-- First 100 lots: 80% Permanent (80), 20% Rent (20) in random order
permanent_rent AS (
    SELECT
        seq::text AS id,
        CASE WHEN row_number() OVER (ORDER BY random()) <= 80
                 THEN 'PERMANENT_PARKING_TYPE'
             ELSE 'RENT_PARKING_TYPE' END AS parking_type
    FROM generate_series(1, 100) AS seq
),
-- Next 10 lots: Special
special AS (
    SELECT generate_series(101, 110)::text AS id,
           'SPECIAL_PARKING_TYPE' AS parking_type
),
-- Last 10 lots: Inclusive
inclusive AS (
    SELECT generate_series(111, 120)::text AS id,
           'INCLUSIVE_PARKING_TYPE' AS parking_type
),
-- Combine all lots
all_lots AS (
    SELECT * FROM permanent_rent
    UNION ALL
    SELECT * FROM special
    UNION ALL
    SELECT * FROM inclusive
),
-- Add random UUID flags (50% chance to populate)
lots_with_uuids AS (
    SELECT *,
           random() < 0.3 AS populate_uuid
    FROM all_lots
)
SELECT
    id,
    parking_type,
    CASE WHEN populate_uuid THEN gen_random_uuid() END AS vehicle_id,
    CASE WHEN populate_uuid THEN gen_random_uuid() END AS owner_id
FROM lots_with_uuids
ORDER BY id::integer;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE parking_lots;
-- +goose StatementEnd
