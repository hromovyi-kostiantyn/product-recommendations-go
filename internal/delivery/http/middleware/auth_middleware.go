// Package middleware implements middleware for the application
package middleware

import (
	"context"
	"net/http"
	"product-recommendations-go/internal/service"
	"strings"
)

// AuthMiddleware реалізує middleware для аутентифікації
type AuthMiddleware struct {
	authService service.AuthService
}

// NewAuthMiddleware створює новий middleware для аутентифікації
func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Middleware виконує перевірку JWT токена
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отримання токена з заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Перевірка формату токена
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Парсинг токена
		userID, err := m.authService.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Додавання ID користувача до контексту запиту
		type contextKey string
		const userIDKey contextKey = "user_id"
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
