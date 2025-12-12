-- Пользователи (администраторы)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Виды спорта
CREATE TABLE sport_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

-- Разряды
CREATE TABLE ranks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    level INTEGER NOT NULL CHECK (level > 0)
);

-- Спортсмены
CREATE TABLE athletes (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),
    birth_date DATE NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
    sport_type_id INTEGER REFERENCES sport_types(id),
    rank_id INTEGER REFERENCES ranks(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Соревнования
CREATE TABLE competitions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    location VARCHAR(200),
    sport_type_id INTEGER REFERENCES sport_types(id)
);

-- Результаты соревнований
CREATE TABLE competition_results (
    id SERIAL PRIMARY KEY,
    competition_id INTEGER REFERENCES competitions(id) ON DELETE CASCADE,
    athlete_id INTEGER REFERENCES athletes(id) ON DELETE CASCADE,
    place INTEGER CHECK (place > 0),
    score DECIMAL(10,2),
    qualification BOOLEAN DEFAULT false,
    UNIQUE(competition_id, athlete_id)
);

-- Индексы для улучшения производительности
CREATE INDEX idx_athletes_sport_type ON athletes(sport_type_id);
CREATE INDEX idx_athletes_rank ON athletes(rank_id);
CREATE INDEX idx_competitions_sport_type ON competitions(sport_type_id);
CREATE INDEX idx_results_competition ON competition_results(competition_id);
CREATE INDEX idx_results_athlete ON competition_results(athlete_id);