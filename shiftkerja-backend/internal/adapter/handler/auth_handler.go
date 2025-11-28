package handler

import (
	"encoding/json"
	"fmt" // 👈 Added this for debugging
	"net/http"
	"shiftkerja-backend/internal/adapter/repository"
	"shiftkerja-backend/internal/core/entity"
	"shiftkerja-backend/internal/core/service"

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
		fmt.Printf("❌ Register JSON Decode Error: %v\n", err) // Debug Log
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 2. Hash the Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("❌ Bcrypt Error: %v\n", err) // Debug Log
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
	fmt.Printf("📝 Attempting to register user: %s\n", user.Email) // Debug Log
	if err := h.Repo.CreateUser(r.Context(), &user); err != nil {
		fmt.Printf("❌ DB CreateUser Error: %v\n", err) // 👈 THIS IS THE IMPORTANT ONE
		// In a real app, check for "duplicate email" error here
		http.Error(w, "Failed to register user (Email might exist)", http.StatusInternalServerError)
		return
	}

	// 5. Respond
	fmt.Printf("✅ User registered successfully: %s\n", user.Email) // Debug Log
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "User created", "email": user.Email})
}

// Request payload for Login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 1. Parse JSON
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("❌ Login JSON Decode Error: %v\n", err) // Debug Log
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Printf("🔑 Attempting login for: %s\n", req.Email) // Debug Log

	// 2. Find User by Email
	user, err := h.Repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		fmt.Printf("❌ DB GetUser Error: %v\n", err) // 👈 Check for DB connection issues
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if user == nil {
		fmt.Printf("❌ User not found in DB: %s\n", req.Email) // 👈 Check if email exists
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 3. Compare Passwords (Hash vs Plain)
	// bcrypt.CompareHashAndPassword(hashed, plain)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		fmt.Printf("❌ Password Mismatch for %s: %v\n", req.Email, err) // 👈 Check if password is wrong
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 4. Generate JWT
	token, err := service.GenerateToken(user.ID, user.Role)
	if err != nil {
		fmt.Printf("❌ Token Generation Error: %v\n", err) // Debug Log
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 5. Respond
	fmt.Printf("✅ Login Successful for %s\n", req.Email) // Debug Log
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"role":  user.Role,
	})
}