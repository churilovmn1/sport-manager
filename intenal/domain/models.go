package domain

import (
    "time"
)

// User - модель пользователя (администратора)
type User struct {
    ID           int       `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    IsAdmin      bool      `json:"is_admin"`
    CreatedAt    time.Time `json:"created_at"`
}

// Athlete - спортсмен
type Athlete struct {
    ID          int       `json:"id"`
    FirstName   string    `json:"first_name"`
    LastName    string    `json:"last_name"`
    MiddleName  string    `json:"middle_name,omitempty"`
    BirthDate   time.Time `json:"birth_date"`
    Gender      string    `json:"gender"`
    SportTypeID int       `json:"sport_type_id"`
    RankID      int       `json:"rank_id"`
    CreatedAt   time.Time `json:"created_at"`
}

// SportType - вид спорта
type SportType struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
}

// Rank - спортивный разряд
type Rank struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Level int    `json:"level"` // 1-высший, 2-средний и т.д.
}

// Competition - соревнование
type Competition struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    Date        time.Time `json:"date"`
    Location    string    `json:"location"`
    SportTypeID int       `json:"sport_type_id"`
}

// CompetitionResult - результат участия в соревновании
type CompetitionResult struct {
    ID             int `json:"id"`
    CompetitionID  int `json:"competition_id"`
    AthleteID      int `json:"athlete_id"`
    Place          int `json:"place"` // 1-первое место и т.д.
    Score          float64 `json:"score,omitempty"`
    Qualification  bool    `json:"qualification"` // прошел квалификацию
}