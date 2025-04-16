package models

// OrderItem представляє товар в замовленні
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"index;not null" json:"order_id"`
	ProductID uint    `gorm:"index;not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"`
	Order     Order   `gorm:"foreignKey:OrderID" json:"-"`
	Product   Product `gorm:"foreignKey:ProductID" json:"-"`
}
