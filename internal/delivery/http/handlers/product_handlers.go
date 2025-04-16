package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"product-recommendations-go/internal/models"
	"product-recommendations-go/internal/service"
	"strconv"
)

// ProductHandler реалізує обробку запитів продуктів
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler створює новий обробник для продуктів
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetAll повертає всі продукти з пагінацією
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Отримуємо параметри пагінації з запиту
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	// Отримуємо продукти з сервісу
	products, total, err := h.productService.GetAll(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Готуємо відповідь
	response := struct {
		Products []*models.Product `json:"products"`
		Total    int64             `json:"total"`
		Page     int               `json:"page"`
		Limit    int               `json:"limit"`
	}{
		Products: products,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	// Повертаємо JSON-відповідь
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error decode JSON", http.StatusInternalServerError)
		log.Printf("Error JSON: %v", err)
		return
	}
}

// GetByID повертає продукт за ID
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID продукту з URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Отримуємо продукт з сервісу
	product, err := h.productService.GetByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Повертаємо JSON-відповідь
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		log.Println("Error JSON encode:", err)
		return
	}
}
