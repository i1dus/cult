-- +goose Up
-- +goose StatementBegin
CREATE TABLE guests
(
    parking_lot_id TEXT NOT NULL,
    start_at       TIMESTAMP,
    end_at         TIMESTAMP,
    created_at     TIMESTAMP DEFAULT now(),
    updated_at     TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE guests;
-- +goose StatementEnd
