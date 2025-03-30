package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"product-recommendations-go/internal/models"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository створює новий екземпляр репозиторію продуктів
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product

	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Продукт не знайдено
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetAll(ctx context.Context, page, limit int) ([]*models.Product, int64, error) {
	var products []*models.Product
	var totalCount int64

	// Рахуємо загальну кількість продуктів
	if err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Якщо сторінка чи ліміт не вказані, встановлюємо значення за замовчуванням
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	// Обчислюємо зміщення
	offset := (page - 1) * limit

	// Отримуємо продукти з пагінацією
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
