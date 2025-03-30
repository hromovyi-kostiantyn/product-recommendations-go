package service

import (
	"context"
	"log"
	"product-recommendations-go/internal/models"
	"product-recommendations-go/internal/repository"
	"product-recommendations-go/pkg/recommendation"
)

type recommendationService struct {
	likeRepo    repository.UserLikeRepository
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

// NewRecommendationService створює новий екземпляр сервісу рекомендацій
func NewRecommendationService(
	likeRepo repository.UserLikeRepository,
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
) RecommendationService {
	return &recommendationService{
		likeRepo:    likeRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *recommendationService) GetRecommendations(ctx context.Context, userID uint, limit int) ([]*models.ProductRecommendation, error) {
	// Якщо ліміт не вказаний або недійсний, встановлюємо значення за замовчуванням
	if limit <= 0 {
		limit = 10
	}

	// Отримуємо лайки користувача
	userLikes, err := s.likeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Отримуємо замовлення користувача
	userOrders, err := s.orderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Отримуємо всі товари
	allProducts, _, err := s.productRepo.GetAll(ctx, 1, 1000) // Отримуємо велику кількість товарів
	if err != nil {
		return nil, err
	}

	log.Printf("User ID: %d, Likes count: %d, Orders count: %d, Products count: %d",
		userID, len(userLikes), len(userOrders), len(allProducts))

	// Викликаємо функцію для обчислення рекомендацій
	recommendedProducts, recommendationScores := recommendation.RecommendProducts(userID, userLikes, userOrders, allProducts, limit)

	log.Printf("Received recommendations: %d, scores: %d", len(recommendedProducts), len(recommendationScores))

	// Переконуємося, що у нас однакова кількість продуктів і оцінок
	if len(recommendedProducts) != len(recommendationScores) {
		log.Printf("Warning: Mismatch between recommendations (%d) and scores (%d)",
			len(recommendedProducts), len(recommendationScores))

		// Якщо є невідповідність, скорочуємо довшу частину
		minLength := len(recommendedProducts)
		if len(recommendationScores) < minLength {
			minLength = len(recommendationScores)
		}

		recommendedProducts = recommendedProducts[:minLength]
		recommendationScores = recommendationScores[:minLength]
	}

	// Підготовка відповіді з рейтингами
	recommendations := make([]*models.ProductRecommendation, 0, len(recommendedProducts))
	for i, product := range recommendedProducts {
		recommendations = append(recommendations, &models.ProductRecommendation{
			Product: product,
			Score:   recommendationScores[i],
		})
	}

	// Обмежуємо кількість рекомендацій
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	log.Printf("Final recommendations with scores: %d", len(recommendations))
	return recommendations, nil
}
