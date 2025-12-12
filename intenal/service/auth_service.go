package service

import (
    "errors"
    "time"
    
    "sport-manager/internal/domain"
    "sport-manager/internal/pkg/jwt"
    "sport-manager/internal/pkg/password"
    "github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
    userRepo domain.UserRepository
    jwt      *jwt.Manager
}

func NewAuthService(userRepo domain.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        jwt:      jwt.NewManager(jwtSecret, 24*time.Hour),
    }
}

func (s *AuthService) Register(username, email, plainPassword string) error {
    hashedPassword, err := password.Hash(plainPassword)
    if err != nil {
        return err
    }
    
    user := &domain.User{
        Username:     username,
        Email:        email,
        PasswordHash: hashedPassword,
        IsAdmin:      true, // регистрируем только админов
        CreatedAt:    time.Now(),
    }
    
    return s.userRepo.Create(user)
}

func (s *AuthService) Login(username, password string) (string, error) {
    user, err := s.userRepo.FindByUsername(username)
    if err != nil {
        return "", errors.New("invalid credentials")
    }
    
    if !password.Verify(password, user.PasswordHash) {
        return "", errors.New("invalid credentials")
    }
    
    token, err := s.jwt.GenerateToken(user.ID, user.IsAdmin)
    if err != nil {
        return "", err
    }
    
    return token, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Claims, error) {
    return s.jwt.ValidateToken(tokenString)
}