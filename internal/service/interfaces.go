package service

import (
	"context"
	"product-recommendations-go/internal/models"
)

// AuthService інтерфейс для роботи з аутентифікацією
type AuthService interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, token string) error
	ParseToken(token string) (uint, error)
}

// ProductService інтерфейс для роботи з товарами
type ProductService interface {
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	GetAll(ctx context.Context, page, limit int) ([]*models.Product, int64, error)
}

// LikeService інтерфейс для роботи з вподобаннями
type LikeService interface {
	LikeProduct(ctx context.Context, userID, productID uint) error
	UnlikeProduct(ctx context.Context, userID, productID uint) error
	GetUserLikes(ctx context.Context, userID uint) ([]*models.UserLike, error)
}

// OrderService інтерфейс для роботи з замовленнями
type OrderService interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, id, userID uint) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID uint) ([]*models.Order, error)
}

// RecommendationService інтерфейс для роботи з рекомендаціями
type RecommendationService interface {
	GetRecommendations(ctx context.Context, userID uint, limit int) ([]*models.ProductRecommendation, error)
}
