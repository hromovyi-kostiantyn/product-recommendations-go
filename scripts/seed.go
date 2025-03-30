package main

import (
	"fmt"
	"log"
	"product-recommendations-go/internal/config"
	"product-recommendations-go/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func main() {
	db := config.GetDB()

	// Створення тестових користувачів
	users := []models.User{
		{Email: "user1@example.com", Password: hashPassword("password1")},
		{Email: "user2@example.com", Password: hashPassword("password2")},
		{Email: "user3@example.com", Password: hashPassword("password3")},
	}

	for _, user := range users {
		db.Create(&user)
	}

	// Створення тестових товарів
	categories := []string{"Electronics", "Clothing", "Books", "Home", "Sports"}
	var products []models.Product

	for i := 1; i <= 50; i++ {
		product := models.Product{
			Name:        fmt.Sprintf("Product %d", i),
			Description: fmt.Sprintf("Description for product %d", i),
			Price:       float64(i*10) + 0.99,
			Category:    categories[i%len(categories)],
			ImageURL:    fmt.Sprintf("https://example.com/images/product%d.jpg", i),
		}
		products = append(products, product)
		db.Create(&product)
	}

	// Створення лайків
	var user1, user2, user3 models.User
	db.First(&user1, 1)
	db.First(&user2, 2)
	db.First(&user3, 3)

	// Користувач 1 лайкає перші 10 товарів
	for i := 1; i <= 10; i++ {
		db.Create(&models.UserLike{
			UserID:    user1.ID,
			ProductID: uint(i),
		})
	}

	// Користувач 2 лайкає товари з 5 по 15
	for i := 5; i <= 15; i++ {
		db.Create(&models.UserLike{
			UserID:    user2.ID,
			ProductID: uint(i),
		})
	}

	// Користувач 3 лайкає товари з 10 по 20
	for i := 10; i <= 20; i++ {
		db.Create(&models.UserLike{
			UserID:    user3.ID,
			ProductID: uint(i),
		})
	}

	// Створення замовлень
	order1 := models.Order{
		UserID:    user1.ID,
		Status:    "completed",
		Total:     159.97,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}
	db.Create(&order1)

	// Додавання товарів у замовлення
	orderItems := []models.OrderItem{
		{OrderID: order1.ID, ProductID: 1, Quantity: 1, Price: 10.99},
		{OrderID: order1.ID, ProductID: 3, Quantity: 2, Price: 30.99},
		{OrderID: order1.ID, ProductID: 5, Quantity: 1, Price: 50.99},
	}
	for _, item := range orderItems {
		db.Create(&item)
	}

	log.Println("Database seeding completed successfully")
}
