




-- Создание таблицы для хранения кодов авторизации
CREATE TABLE IF NOT EXISTS auth_codes (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    used BOOLEAN DEFAULT FALSE
);