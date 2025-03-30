package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"product-recommendations-go/internal/models"
)

type userLikeRepository struct {
	db *gorm.DB
}

// NewUserLikeRepository створює новий екземпляр репозиторію лайків
func NewUserLikeRepository(db *gorm.DB) UserLikeRepository {
	return &userLikeRepository{
		db: db,
	}
}

func (r *userLikeRepository) Create(ctx context.Context, like *models.UserLike) error {
	return r.db.WithContext(ctx).Create(like).Error
}

func (r *userLikeRepository) Delete(ctx context.Context, userID, productID uint) error {
	result := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.UserLike{})
	if result.Error != nil {
		return result.Error
	}

	// Перевіряємо, чи були видалені записи
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}

func (r *userLikeRepository) GetByUserID(ctx context.Context, userID uint) ([]*models.UserLike, error) {
	var likes []*models.UserLike

	// Завантажуємо також інформацію про продукти
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Where("user_id = ?", userID).
		Find(&likes).Error; err != nil {
		return nil, err
	}

	return likes, nil
}

func (r *userLikeRepository) Exists(ctx context.Context, userID, productID uint) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserLike{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
