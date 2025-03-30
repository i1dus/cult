-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
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

INSERT INTO users (id, phone, name, surname, patronymic, address, user_type, pass_hash)
VALUES
    ('4f8d4d2c-719c-4fad-a7d5-27cd5b515a56', '+79001234567', 'Иван', 'Иванов', 'Иванович', 'ул. Пушкина, д.10', 'REGULAR_USER_TYPE', '$2a$10$7N8bQ4H2jklVdDfDp/YqE.9r4Jp7W1T0KjJv8mR6Y1Z3sLbNcXrO'),
    ('4d0cc3c1-7ad4-4a16-8120-baf081845fae', '+79007654321', 'Мария', 'Петрова', 'Сергеевна', 'пр. Ленина, д.25', 'ADMINISTRATOR_USER_TYPE', '$2a$10$hH5sTqW3vDdQkLpM9oRt0uV1AeR2fS3gU4iY5jK6l7mN8bV7cX9z'),
    ('094d3560-e1c8-4a71-a786-9b94243e8d18', '+79005553535', 'Алексей', 'Смирнов', 'Андреевич', 'ул. Гагарина, д.5', 'MANAGING_COMPANY_USER_TYPE', '$2a$10$3kLmNoPqRsTuVwXyZzAbC.1r2Jp7W1T0KjJv8mR6Y1Z3sLbNcXrO'),
    ('0ae69239-ac40-4046-b93c-89b098b5432f', '+79008887766', 'Ольга', 'Соколова', 'Дмитриевна', 'пр. Мира, д.15', 'REGULAR_USER_TYPE', '$2a$10$9QwErTyUiOpAsDfGhJkLz.4r5Jp7W1T0KjJv8mR6Y1Z3sLbNcXrO'),
    ('b854f850-9048-4b7f-b531-946f3f61b84b','+79002345678', 'Дмитрий', 'Кузнецов', NULL, 'ул. Садовая, д.3', 'REGULAR_USER_TYPE', '$2a$10$2WsXrCtvGyBnUmH4jKl.O5t6Jp7W1T0KjJv8mR6Y1Z3sLbNcXrO'),
    ('642603ee-777d-4264-8270-4edd2f61efab','+79008765432', 'Екатерина', 'Попова', 'Алексеевна', NULL, 'REGULAR_USER_TYPE', '$2a$10$5tGyHujMkOlPqRsTuVwXyZ'),
    ('47cc19dd-6411-43ab-a929-e15ded662c93','+79004567890', 'Сергей', 'Васильев', 'Петрович', 'ул. Лесная, д.7', 'MANAGING_COMPANY_USER_TYPE', '$2a$10$7H8jK9L0Z1X2C3V4B5N6M'),
    ('2cc6d4ff-c7c0-4f19-8a2c-357450296fb1','+79003215647', 'Анна', 'Павлова', 'Викторовна', 'пр. Космонавтов, д.12', 'REGULAR_USER_TYPE', '$2a$10$1Q2W3E4R5T6Y7U8I9O0P'),
    ('e19c9873-df4f-4421-9d2f-20efcc4b2c18','+79006543218', 'Павел', 'Семенов', 'Олегович', NULL, 'ADMINISTRATOR_USER_TYPE', '$2a$10$6Y5U4I3O2P1A9S8D7F6G'),
    ('70d89566-dddf-4f76-b28a-f5e03fdc87d8','+79007778899', 'Татьяна', 'Федорова', 'Николаевна', 'ул. Центральная, д.1', 'REGULAR_USER_TYPE', '$2a$10$8H7J6G5F4D3S2A1Q2W3E'),
    ('bbf1f5f1-d0fe-4946-bd71-90b15cd40a17','+79007778810', NULL, NULL, NULL, NULL, 'REGULAR_USER_TYPE', '$2a$10$8H7J6G5F4D3S2A1Q2W3E'),
    ('9100aeb3-ac9a-4941-9c86-4cbee5da8b71','+79000000000', 'Админ', 'Админов', 'Админович', 'ул. Пушкина, д.1', 'ADMINISTRATOR_USER_TYPE', '$2a$10$nmONwfFLIAyQMPXbk618tuFwrvbQGw1U590dqWcVvCHhGKgvAyjxm'),
    ('3374d806-a42e-42ff-9fbe-b7ea6e29efc0','+79000000001', 'Компания', 'Управляющая', NULL, 'ул. Гоголя, д.1', 'MANAGING_COMPANY_USER_TYPE', '$2a$10$nmONwfFLIAyQMPXbk618tuFwrvbQGw1U590dqWcVvCHhGKgvAyjxm');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
