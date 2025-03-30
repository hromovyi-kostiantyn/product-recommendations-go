package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"product-recommendations-go/internal/models"
)

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository створює новий екземпляр репозиторію замовлень
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order

	if err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Замовлення не знайдено
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uint) ([]*models.Order, error) {
	var orders []*models.Order

	if err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Where("user_id = ?", userID).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) AddItem(ctx context.Context, orderItem *models.OrderItem) error {
	return r.db.WithContext(ctx).Create(orderItem).Error
}
