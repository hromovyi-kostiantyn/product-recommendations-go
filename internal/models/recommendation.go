package models

type ProductRecommendation struct {
	Product *Product `json:"product"`
	Score   float64  `json:"score"`
}
