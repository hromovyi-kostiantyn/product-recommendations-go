// Package models provides the data models for the application.
package models

import (
	"gorm.io/gorm"
	"time"
)

// Order represents a user's order in the system.
type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Status    string         `gorm:"not null" json:"status"`
	Total     float64        `json:"total"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Items     []OrderItem    `json:"items"`
}
