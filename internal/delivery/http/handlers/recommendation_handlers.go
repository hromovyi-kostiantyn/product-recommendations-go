package handlers

import (
	"encoding/json"
	"net/http"
	"product-recommendations-go/internal/models"
	"product-recommendations-go/internal/service"
	"strconv"
)

type RecommendationHandler struct {
	recommendationService service.RecommendationService
}

// NewRecommendationHandler створює новий обробник для рекомендацій
func NewRecommendationHandler(recommendationService service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationService: recommendationService,
	}
}

// GetRecommendations повертає рекомендації продуктів для користувача
func (h *RecommendationHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	recommendations, err := h.recommendationService.GetRecommendations(r.Context(), userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Переконуємося, що повертаємо порожній масив, а не null
	if recommendations == nil {
		recommendations = []*models.ProductRecommendation{}
	}

	// Повертаємо JSON-відповідь
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Recommendations []*models.ProductRecommendation `json:"recommendations"`
	}{
		Recommendations: recommendations,
	})
}
