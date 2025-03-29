-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exceptions (

    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
