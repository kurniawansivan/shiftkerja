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

// --- GET NEARBY (Existing) ---
func (h *ShiftHandler) GetNearby(w http.ResponseWriter, r *http.Request) {
	// ... (Keep your existing logic here, or see below if you want the clean version) ...
	// For brevity, assuming you keep the existing GetNearby logic. 
	// Make sure to parse lat/lng/rad and call h.RedisRepo.FindNearby
	// ---------------------------------------------------------
	
	// Quick re-implementation just in case:
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	lat, _ := strconv.ParseFloat(q.Get("lat"), 64)
	lng, _ := strconv.ParseFloat(q.Get("lng"), 64)
	rad, _ := strconv.ParseFloat(q.Get("rad"), 64)
	if rad == 0 { rad = 10 }

	shifts, err := h.RedisRepo.FindNearby(r.Context(), lat, lng, rad)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(shifts)
}

// --- CREATE SHIFT (New) ---
func (h *ShiftHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Security Check: Are you a business?
	// The middleware put "role" into the context
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
	// Important: We use the Postgres ID we just generated
	// Convert int64 ID to string for Redis
	req.ID = req.ID // It's already set by Scan() in PostgresRepo
	
	// Note: We might need to adjust RedisRepo.AddShift because ID is now int64 in struct
	// See Step 3.5 below
	if err := h.RedisRepo.AddShift(r.Context(), req); err != nil {
		// Log error but don't fail request (Data is safe in SQL)
		fmt.Printf("⚠️ Redis Sync Error: %v\n", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}