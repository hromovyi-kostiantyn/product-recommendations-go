package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"product-recommendations-go/internal/service"
	"strconv"
)

// LikeHandler реалізує обробку запитів лайків
type LikeHandler struct {
	likeService service.LikeService
}

// NewLikeHandler створює новий обробник для лайків
func NewLikeHandler(likeService service.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

// LikeProduct додає продукт до лайків користувача
func (h *LikeHandler) LikeProduct(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID продукту з URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["product_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Отримуємо ID користувача з контексту (встановлений middleware аутентифікації)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Додаємо лайк
	err = h.likeService.LikeProduct(r.Context(), userID, uint(productID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Повертаємо успішну відповідь
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(`{"message":"Product liked successfully"}`))
	if err != nil {
		log.Printf("Error in response: %v", err)
	}
}

// UnlikeProduct видаляє продукт з лайків користувача
func (h *LikeHandler) UnlikeProduct(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID продукту з URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["product_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Отримуємо ID користувача з контексту (встановлений middleware аутентифікації)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Видаляємо лайк
	err = h.likeService.UnlikeProduct(r.Context(), userID, uint(productID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Повертаємо успішну відповідь
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message":"Product unliked successfully"}`))
	if err != nil {
		log.Printf("Error in response: %v", err)
	}
}

// GetUserLikes повертає всі лайкнуті продукти користувача
func (h *LikeHandler) GetUserLikes(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID користувача з контексту (встановлений middleware аутентифікації)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Отримуємо лайки користувача
	likes, err := h.likeService.GetUserLikes(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Повертаємо JSON-відповідь
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Likes interface{} `json:"likes"`
	}{
		Likes: likes,
	})
	if err != nil {
		log.Printf("Error JSON: %v", err)
	}
}
