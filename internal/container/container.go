// Package container DI dependency injection
package container

import (
	"product-recommendations-go/internal/config"
	"product-recommendations-go/internal/delivery/http/handlers"
	"product-recommendations-go/internal/repository"
	"product-recommendations-go/internal/service"
)

// Container зберігає всі залежності програми
type Container struct {
	// Репозиторії
	UserRepository    repository.UserRepository
	ProductRepository repository.ProductRepository
	LikeRepository    repository.UserLikeRepository
	OrderRepository   repository.OrderRepository

	// Сервіси
	AuthService           service.AuthService
	ProductService        service.ProductService
	LikeService           service.LikeService
	OrderService          service.OrderService
	RecommendationService service.RecommendationService

	// Обробники HTTP запитів
	AuthHandler           *handlers.AuthHandler
	ProductHandler        *handlers.ProductHandler
	LikeHandler           *handlers.LikeHandler
	OrderHandler          *handlers.OrderHandler
	RecommendationHandler *handlers.RecommendationHandler
}

// NewContainer створює новий контейнер залежностей
func NewContainer() *Container {
	// Отримуємо підключення до бази даних
	db := config.GetDB()

	// Ініціалізуємо репозиторії
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	likeRepo := repository.NewUserLikeRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Отримуємо JWT секретний ключ
	jwtSecret := config.GetEnv("JWT_SECRET", "your-secret-key")

	// Ініціалізуємо сервіси
	authService := service.NewAuthService(userRepo, jwtSecret)
	productService := service.NewProductService(productRepo)
	likeService := service.NewLikeService(likeRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo)
	recommendationService := service.NewRecommendationService(likeRepo, orderRepo, productRepo)

	// Ініціалізуємо обробники
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	likeHandler := handlers.NewLikeHandler(likeService)
	orderHandler := handlers.NewOrderHandler(orderService)
	recommendationHandler := handlers.NewRecommendationHandler(recommendationService)

	// Створюємо контейнер
	return &Container{
		UserRepository:    userRepo,
		ProductRepository: productRepo,
		LikeRepository:    likeRepo,
		OrderRepository:   orderRepo,

		AuthService:           authService,
		ProductService:        productService,
		LikeService:           likeService,
		OrderService:          orderService,
		RecommendationService: recommendationService,

		AuthHandler:           authHandler,
		ProductHandler:        productHandler,
		LikeHandler:           likeHandler,
		OrderHandler:          orderHandler,
		RecommendationHandler: recommendationHandler,
	}
}
