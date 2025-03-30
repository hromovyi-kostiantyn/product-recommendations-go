package service

import (
	"context"
	"product-recommendations-go/internal/models"
	"product-recommendations-go/internal/repository"
)

type productService struct {
	productRepo repository.ProductRepository
}

// NewProductService створює новий екземпляр сервісу продуктів
func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) GetAll(ctx context.Context, page, limit int) ([]*models.Product, int64, error) {
	return s.productRepo.GetAll(ctx, page, limit)
}
