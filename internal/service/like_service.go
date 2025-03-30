package service

import (
	"context"
	"errors"
	"product-recommendations-go/internal/models"
	"product-recommendations-go/internal/repository"
)

type likeService struct {
	likeRepo    repository.UserLikeRepository
	productRepo repository.ProductRepository
}

// NewLikeService створює новий екземпляр сервісу лайків
func NewLikeService(likeRepo repository.UserLikeRepository, productRepo repository.ProductRepository) LikeService {
	return &likeService{
		likeRepo:    likeRepo,
		productRepo: productRepo,
	}
}

func (s *likeService) LikeProduct(ctx context.Context, userID, productID uint) error {
	// Перевіряємо, чи існує продукт
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	if product == nil {
		return errors.New("product not found")
	}

	// Перевіряємо, чи користувач уже лайкав цей продукт
	exists, err := s.likeRepo.Exists(ctx, userID, productID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("user already liked this product")
	}

	// Створюємо новий лайк
	like := &models.UserLike{
		UserID:    userID,
		ProductID: productID,
	}

	return s.likeRepo.Create(ctx, like)
}

func (s *likeService) UnlikeProduct(ctx context.Context, userID, productID uint) error {
	// Видаляємо лайк
	return s.likeRepo.Delete(ctx, userID, productID)
}

func (s *likeService) GetUserLikes(ctx context.Context, userID uint) ([]*models.UserLike, error) {
	return s.likeRepo.GetByUserID(ctx, userID)
}
