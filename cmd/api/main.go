// Package main entry point for the application
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"product-recommendations-go/internal/config"
	"product-recommendations-go/internal/container"
	"product-recommendations-go/internal/delivery/http/middleware"
)

func main() {
	// Створюємо контейнер залежностей
	c := container.NewContainer()

	// Створюємо маршрутизатор
	r := mux.NewRouter()

	// Маршрути для аутентифікації (публічні)
	r.HandleFunc("/api/v1/auth/register", c.AuthHandler.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", c.AuthHandler.Login).Methods("POST")

	// Middleware для перевірки JWT токена
	authMiddleware := middleware.NewAuthMiddleware(c.AuthService)

	// Захищені маршрути (потрібна аутентифікація)
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(authMiddleware.Middleware)

	// Маршрут для виходу з системи
	api.HandleFunc("/auth/logout", c.AuthHandler.Logout).Methods("POST")

	// Маршрути для товарів
	api.HandleFunc("/products", c.ProductHandler.GetAll).Methods("GET")
	api.HandleFunc("/products/{id}", c.ProductHandler.GetByID).Methods("GET")

	// Маршрути для лайків
	api.HandleFunc("/likes/{product_id}", c.LikeHandler.LikeProduct).Methods("POST")
	api.HandleFunc("/likes/{product_id}", c.LikeHandler.UnlikeProduct).Methods("DELETE")
	api.HandleFunc("/likes", c.LikeHandler.GetUserLikes).Methods("GET")

	// Маршрути для замовлень
	api.HandleFunc("/orders", c.OrderHandler.CreateOrder).Methods("POST")
	api.HandleFunc("/orders", c.OrderHandler.GetUserOrders).Methods("GET")
	api.HandleFunc("/orders/{id}", c.OrderHandler.GetOrderByID).Methods("GET")

	// Маршрути для рекомендацій
	api.HandleFunc("/recommendations", c.RecommendationHandler.GetRecommendations).Methods("GET")

	// Перевірка стану сервісу (без аутентифікації)
	r.HandleFunc("/api/v1/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			log.Println("Error writing response:", err)
			return
		}
	}).Methods("GET")

	// Налаштування сервера
	port := config.GetEnv("APP_PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в окремій горутині
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Перехоплення сигналів для граційного завершення
	sigChan := make(chan os.Signal, 1) // Змінили ім'я змінної
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Очікування сигналу
	<-sigChan

	// Створення контексту з таймаутом для завершення обробки запитів
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Завершення роботи сервера
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	// Закриття підключення до бази даних
	config.CloseDB()

	log.Println("Server stopped gracefully")
}
