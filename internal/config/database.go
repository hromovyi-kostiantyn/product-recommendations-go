// Package config represent the configuration settings for the application
package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
	mu       sync.Mutex
)

// GetDB повертає єдиний екземпляр підключення до бази даних
func GetDB() *gorm.DB {
	once.Do(func() {
		dbHost := GetEnv("DB_HOST", "localhost")
		dbUser := GetEnv("DB_USER", "postgres")
		dbPassword := GetEnv("DB_PASSWORD", "postgres")
		dbName := GetEnv("DB_NAME", "recommendations")
		dbPort := GetEnv("DB_PORT", "5432")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbHost, dbUser, dbPassword, dbName, dbPort)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Встановлення налаштувань пулу підключень
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get DB instance: %v", err)
		}

		// Встановлюємо ліміти підключень
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)

		instance = db
		log.Println("Database connection established")
	})

	return instance
}

// CloseDB закриває підключення до бази даних
func CloseDB() {
	mu.Lock()
	defer mu.Unlock()
	if instance != nil {
		sqlDB, err := instance.DB()
		if err != nil {
			log.Printf("Error getting DB instance: %v", err)
			return
		}

		if err := sqlDB.Close(); err != nil {
			log.Printf("BD connection error: %v", err)
		}
		instance = nil
		log.Println("Database connection closed")
	}
}

// GetEnv отримує значення з змінних середовища або повертає запасне значення
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
