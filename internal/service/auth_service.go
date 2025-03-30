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

func (s *authService) Logout(ctx context.Context, token string) error {
	// Для простого JWT досить не зберігати токен на стороні клієнта
	// В реальному проєкті можна додати токен до чорного списку
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
