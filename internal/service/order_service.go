package service

import (
	"context"
	"errors"
	"product-recommendations-go/internal/models"
	"product-recommendations-go/internal/repository"
)

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

// NewOrderService створює новий екземпляр сервісу замовлень
func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order *models.Order) error {
	// Перевіряємо, чи є товари в замовленні
	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	// Обчислюємо суму замовлення та перевіряємо товари
	var total float64
	for i, item := range order.Items {
		// Отримуємо товар для перевірки його наявності та поточної ціни
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return err
		}

		if product == nil {
			return errors.New("product not found")
		}

		// Встановлюємо поточну ціну
		order.Items[i].Price = product.Price

		// Обчислюємо суму для товару
		itemTotal := product.Price * float64(item.Quantity)
		total += itemTotal
	}

	// Встановлюємо загальну суму замовлення
	order.Total = total

	// Створюємо замовлення
	return s.orderRepo.Create(ctx, order)
}

func (s *orderService) GetOrderByID(ctx context.Context, id, userID uint) (*models.Order, error) {
	// Отримуємо замовлення
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	// Перевіряємо, чи замовлення належить користувачу
	if order.UserID != userID {
		return nil, errors.New("unauthorized access to order")
	}

	return order, nil
}

func (s *orderService) GetUserOrders(ctx context.Context, userID uint) ([]*models.Order, error) {
	return s.orderRepo.GetByUserID(ctx, userID)
}
