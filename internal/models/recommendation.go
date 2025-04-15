package models

// ProductRecommendation представляє рекомендацію продукту
type ProductRecommendation struct {
	Product *Product `json:"product"`
	Score   float64  `json:"score"`
}
