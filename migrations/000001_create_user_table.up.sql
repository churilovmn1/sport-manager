-- Таблица для Администраторов
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL, -- Здесь будем хранить хеш пароля
    is_admin BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Добавление начального пользователя-администратора (пароль: admin123)
-- Настоящий хеш будет создан в коде Go, пока используем заглушку.
INSERT INTO users (username, password_hash)
VALUES ('admin', '$2a$10$sjFtm14QpMa8F8K3aOt0veJsdCVSfzJG1oy1z9VanAU0JDFxeSSIi');