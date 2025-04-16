// Package handlers implements HTTP handlers for the application
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"product-recommendations-go/internal/models"
	"product-recommendations-go/internal/service"
)

// AuthHandler реалізує обробку запитів аутентифікації
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler створює новий обробник для аутентифікації
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

// Register обробляє запит на реєстрацію нового користувача
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.authService.Register(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(`{"message":"User registered successfully"}`))
	if err != nil {
		log.Printf("Error in response: %v", err)
	}
}

// Login обробляє запит на вхід користувача
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := authResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error JSON decode", http.StatusInternalServerError)
		log.Printf("Error JSON: %v", err)
		return
	}
}

// Logout обробляє запит на вихід користувача
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Отримання токена з заголовка Authorization
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header is required", http.StatusBadRequest)
		return
	}

	// Видалення префікса "Bearer "
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := h.authService.Logout(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message":"Logged out successfully"}`))
	if err != nil {
		log.Printf("Error in response: %v", err)
	}
}
