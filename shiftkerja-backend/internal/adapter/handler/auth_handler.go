package handler

import (
	"encoding/json"
	"net/http"
	"shiftkerja-backend/internal/adapter/repository"
	"shiftkerja-backend/internal/core/entity"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Repo *repository.PostgresUserRepo
}

func NewAuthHandler(repo *repository.PostgresUserRepo) *AuthHandler {
	return &AuthHandler{Repo: repo}
}

// Request payload struct
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// 1. Parse JSON Request
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 2. Hash the Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// 3. Create User Entity
	user := entity.User{
		Email:        req.Email,
		PasswordHash: string(hashed),
		FullName:     req.FullName,
		Role:         req.Role,
	}

	// 4. Save to DB
	if err := h.Repo.CreateUser(r.Context(), &user); err != nil {
		// In a real app, check for "duplicate email" error here
		http.Error(w, "Failed to register user (Email might exist)", http.StatusInternalServerError)
		return
	}

	// 5. Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "User created", "email": user.Email})
}