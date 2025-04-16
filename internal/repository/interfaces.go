// Package repository реалізація інтерфейсів для роботи з базою даних
package repository

import (
	"context"
	"product-recommendations-go/internal/models"
)

// UserRepository інтерфейс для роботи з користувачами
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

// ProductRepository інтерфейс для роботи з товарами
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	GetAll(ctx context.Context, page, limit int) ([]*models.Product, int64, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
}

// UserLikeRepository інтерфейс для роботи з вподобаннями
type UserLikeRepository interface {
	Create(ctx context.Context, like *models.UserLike) error
	Delete(ctx context.Context, userID, productID uint) error
	GetByUserID(ctx context.Context, userID uint) ([]*models.UserLike, error)
	Exists(ctx context.Context, userID, productID uint) (bool, error)
}

// OrderRepository інтерфейс для роботи з замовленнями
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uint) (*models.Order, error)
	GetByUserID(ctx context.Context, userID uint) ([]*models.Order, error)
	AddItem(ctx context.Context, orderItem *models.OrderItem) error
}
