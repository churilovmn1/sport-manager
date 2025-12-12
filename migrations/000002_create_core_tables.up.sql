-- Таблица-справочник: Виды спорта
CREATE TABLE IF NOT EXISTS sports (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Таблица-справочник: Разряды
CREATE TABLE IF NOT EXISTS ranks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
);

-- Таблица-справочник: Соревнования
CREATE TABLE IF NOT EXISTS competitions (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    organizer VARCHAR(150),
    UNIQUE (title, start_date)
);

-- Таблица: Спортсмен (основной объект работы)
CREATE TABLE IF NOT EXISTS athletes (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE,
    gender VARCHAR(10), -- 'Male', 'Female', 'Other'
    sport_id INT REFERENCES sports(id) ON DELETE SET NULL, -- Связь с Видом спорта
    rank_id INT REFERENCES ranks(id) ON DELETE SET NULL, -- Связь с Разрядом
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (full_name, birth_date)
);