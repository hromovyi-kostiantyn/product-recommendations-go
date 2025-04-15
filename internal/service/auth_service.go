// Package service надає реалізацію сервісу аутентифікації
package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"product-recommendations-go/internal/models"
	"product-recommendations-go/internal/repository"
	"time"
)

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService створює новий екземпляр сервісу аутентифікації
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, user *models.User) error {
	// Перевірка, чи існує користувач з такою адресою
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Хешування пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Створення користувача
	return s.userRepo.Create(ctx, user)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// Отримання користувача за email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Перевірка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Створення JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Підписання токена
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) Logout(_ context.Context, token string) error {
	// Parse the token to extract expiration time
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return err
	}

	// Ensure the token is valid
	if !parsedToken.Valid {
		return errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Get expiration time
	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid exp claim")
	}

	expiresAt := time.Unix(int64(exp), 0)

	// For a complete implementation, you would:
	// 1. Add this token to a blacklist storage (Redis/DB)
	// 2. Set the blacklist entry to expire at the token's expiration time
	// 3. Check the blacklist in ParseToken before validating tokens

	// Placeholder return until blacklist storage is implemented
	_ = expiresAt // Avoid unused variable warning
	return nil
}

func (s *authService) ParseToken(token string) (uint, error) {
	// Парсинг JWT токена
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Перевірка методу підпису
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	// Перевірка валідності токена
	if !parsedToken.Valid {
		return 0, errors.New("invalid token")
	}

	// Отримання claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// Отримання ID користувача
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id claim")
	}

	return uint(userID), nil
}
