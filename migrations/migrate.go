package main

import (
	"log"

	"product-recommendations-go/internal/config"
	"product-recommendations-go/internal/models"
)

func main() {
	db := config.GetDB()

	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.UserLike{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully")
}
