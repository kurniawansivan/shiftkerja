package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"shiftkerja-backend/internal/adapter/repository"
	"shiftkerja-backend/internal/core/entity"
)

type ShiftHandler struct {
	RedisRepo    *repository.RedisGeoRepository
	PostgresRepo *repository.PostgresShiftRepo
}

// Inject BOTH repos
func NewShiftHandler(redis *repository.RedisGeoRepository, pg *repository.PostgresShiftRepo) *ShiftHandler {
	return &ShiftHandler{
		RedisRepo:    redis,
		PostgresRepo: pg,
	}
}

// --- GET NEARBY (For Map) ---
func (h *ShiftHandler) GetNearby(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Parse Query Params
	q := r.URL.Query()
	lat, _ := strconv.ParseFloat(q.Get("lat"), 64)
	lng, _ := strconv.ParseFloat(q.Get("lng"), 64)
	rad, _ := strconv.ParseFloat(q.Get("rad"), 64)
	if rad == 0 {
		rad = 10 // Default 10km radius
	}

	// Call Redis Repo
	shifts, err := h.RedisRepo.FindNearby(r.Context(), lat, lng, rad)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(shifts)
}

// --- CREATE SHIFT (Business Only) ---
func (h *ShiftHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Security Check: Are you a business?
	role := r.Context().Value("role").(string)
	userID := r.Context().Value("user_id").(float64) // JWT numbers are float64 by default

	if role != "business" && role != "admin" {
		http.Error(w, "Only businesses can post shifts", http.StatusForbidden)
		return
	}

	// 2. Parse Input
	var req entity.Shift
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set the owner from the Token (don't trust the body!)
	req.OwnerID = int64(userID)

	// 3. Save to Postgres (The Source of Truth)
	if err := h.PostgresRepo.CreateShift(r.Context(), &req); err != nil {
		fmt.Printf("❌ PG Error: %v\n", err)
		http.Error(w, "Failed to save shift", http.StatusInternalServerError)
		return
	}

	// 4. Save to Redis (The Geo Index)
	// The req.ID is now populated from Postgres
	if err := h.RedisRepo.AddShift(r.Context(), req); err != nil {
		// Log error but don't fail request (Data is safe in SQL)
		fmt.Printf("⚠️ Redis Sync Error: %v\n", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// --- APPLY FOR SHIFT (Worker Only) ---
// Define Request Body Struct
type ApplyRequest struct {
	ShiftID int64 `json:"shift_id"`
}

func (h *ShiftHandler) Apply(w http.ResponseWriter, r *http.Request) {
	// 1. Who is the user?
	// Note: user_id comes from JWT middleware as float64
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "worker" {
		http.Error(w, "Only workers can apply", http.StatusForbidden)
		return
	}

	// 2. Which shift?
	var req ApplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 3. Save to DB
	err := h.PostgresRepo.ApplyForShift(r.Context(), req.ShiftID, userID)
	if err != nil {
		// Log the specific error for debugging
		fmt.Printf("❌ Apply Error: %v\n", err)
		http.Error(w, "Failed to apply (Shift might not exist or already applied)", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Applied successfully"})
}