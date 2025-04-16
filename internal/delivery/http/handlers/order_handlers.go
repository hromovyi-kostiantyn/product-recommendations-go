package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"product-recommendations-go/internal/models"
	"product-recommendations-go/internal/service"
	"strconv"
)

// OrderHandler реалізує обробку запитів замовлень
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler створює новий обробник для замовлень
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder обробляє створення нового замовлення
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID користувача з контексту (встановлений middleware аутентифікації)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Розбираємо тіло запиту
	var requestOrder struct {
		Items []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestOrder); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Створюємо об'єкт замовлення
	order := &models.Order{
		UserID: userID,
		Status: "pending",
		Items:  make([]models.OrderItem, len(requestOrder.Items)),
	}

	// Додаємо товари до замовлення
	for i, item := range requestOrder.Items {
		order.Items[i] = models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	// Створюємо замовлення
	if err := h.orderService.CreateOrder(r.Context(), order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Повертаємо успішну відповідь
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(struct {
		OrderID uint   `json:"order_id"`
		Message string `json:"message"`
	}{
		OrderID: order.ID,
		Message: "Order created successfully",
	})
	if err != nil {
		http.Error(w, "Error JSON encode", http.StatusInternalServerError)
		return
	}
}

// GetOrderByID повертає замовлення за ID
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID користувача з контексту (встановлений middleware аутентифікації)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Отримуємо ID замовлення з URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Отримуємо замовлення з сервісу
	order, err := h.orderService.GetOrderByID(r.Context(), uint(id), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Повертаємо JSON-відповідь
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		http.Error(w, "Error JSON encode", http.StatusInternalServerError)
		return
	}
}

// GetUserOrders повертає всі замовлення користувача
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	// Отримуємо ID користувача з контексту (встановлений middleware аутентифікації)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Отримуємо замовлення з сервісу
	orders, err := h.orderService.GetUserOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Повертаємо JSON-відповідь
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Orders []*models.Order `json:"orders"`
	}{
		Orders: orders,
	})
	if err != nil {
		http.Error(w, "Error JSON encode", http.StatusInternalServerError)
		return
	}
}
