-- Таблица: Участие в соревновании
CREATE TABLE IF NOT EXISTS participations (
    id SERIAL PRIMARY KEY,
    athlete_id INT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
    competition_id INT NOT NULL REFERENCES competitions(id) ON DELETE CASCADE,
    registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Уникальный индекс: спортсмен может участвовать в одном соревновании только один раз
    UNIQUE (athlete_id, competition_id)
);

-- Таблица: Результаты (место)
CREATE TABLE IF NOT EXISTS results (
    id SERIAL PRIMARY KEY,
    participation_id INT NOT NULL UNIQUE REFERENCES participations(id) ON DELETE CASCADE,
    place INT, -- Место (например, 1, 2, 3)
    score NUMERIC(10, 2), -- Результат (время, очки и т.д.)
    notes TEXT,
    
    -- Проверка на валидность места (не может быть 0 или отрицательным)
    CHECK (place > 0)
);