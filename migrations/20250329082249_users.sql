-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone      TEXT NOT NULL UNIQUE,
    name       TEXT,
    surname    TEXT,
    patronymic TEXT,
    address    TEXT,
    user_type  TEXT NOT NULL,
    pass_hash  TEXT NOT NULL
);

INSERT INTO users (phone, name, surname, patronymic, address, user_type, pass_hash)
VALUES
    ('+79001234567', 'Иван', 'Иванов', 'Иванович', 'ул. Пушкина, д.10', 'REGULAR_USER_TYPE', '$2a$10$7N8bQ4H2jklVdDfDp/YqE.9r4Jp7W1T0KjJv8mR6Y1Z3sLbNcXrO'),
    ('+79007654321', 'Мария', 'Петрова', 'Сергеевна', 'пр. Ленина, д.25', 'ADMINISTRATOR_USER_TYPE', '$2a$10$hH5sTqW3vDdQkLpM9oRt0uV1AeR2fS3gU4iY5jK6l7mN8bV7cX9z'),
    ('+79005553535', 'Алексей', 'Смирнов', 'Андреевич', 'ул. Гагарина, д.5', 'MANAGING_COMPANY_USER_TYPE', '$2a$10$3kLmNoPqRsTuVwXyZzAbC.1r2Jp7W1T0KjJv8mR6Y1Z3sLbNcXrO'),
    ('+79008887766', 'Ольга', 'Соколова', 'Дмитриевна', 'пр. Мира, д.15', 'REGULAR_USER_TYPE', '$2a$10$9QwErTyUiOpAsDfGhJkLz.4r5Jp7W1T0KjJv8mR6Y1Z3sLbNcXrO'),
    ('+79002345678', 'Дмитрий', 'Кузнецов', NULL, 'ул. Садовая, д.3', 'REGULAR_USER_TYPE', '$2a$10$2WsXrCtvGyBnUmH4jKl.O5t6Jp7W1T0KjJv8mR6Y1Z3sLbNcXrO'),
    ('+79008765432', 'Екатерина', 'Попова', 'Алексеевна', NULL, 'REGULAR_USER_TYPE', '$2a$10$5tGyHujMkOlPqRsTuVwXyZ'),
    ('+79004567890', 'Сергей', 'Васильев', 'Петрович', 'ул. Лесная, д.7', 'MANAGING_COMPANY_USER_TYPE', '$2a$10$7H8jK9L0Z1X2C3V4B5N6M'),
    ('+79003215647', 'Анна', 'Павлова', 'Викторовна', 'пр. Космонавтов, д.12', 'REGULAR_USER_TYPE', '$2a$10$1Q2W3E4R5T6Y7U8I9O0P'),
    ('+79006543218', 'Павел', 'Семенов', 'Олегович', NULL, 'ADMINISTRATOR_USER_TYPE', '$2a$10$6Y5U4I3O2P1A9S8D7F6G'),
    ('+79007778899', 'Татьяна', 'Федорова', 'Николаевна', 'ул. Центральная, д.1', 'REGULAR_USER_TYPE', '$2a$10$8H7J6G5F4D3S2A1Q2W3E');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
