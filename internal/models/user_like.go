package models

import "time"

type UserLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index:idx_user_product,unique;not null" json:"user_id"`
	ProductID uint      `gorm:"index:idx_user_product,unique;not null" json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"-"`
}
