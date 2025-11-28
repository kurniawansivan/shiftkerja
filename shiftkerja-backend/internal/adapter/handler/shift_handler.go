package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"shiftkerja-backend/internal/adapter/repository"
)

type ShiftHandler struct {
	Repo *repository.RedisGeoRepository
}

// Dependency Injection: We pass the repo in when we start
func NewShiftHandler(repo *repository.RedisGeoRepository) *ShiftHandler {
	return &ShiftHandler{Repo: repo}
}

func (h *ShiftHandler) GetNearby(w http.ResponseWriter, r *http.Request) {
	// 1. Enable CORS (So Vue running on port 5173 can talk to Go on 8080)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// 2. Parse Query Parameters (lat, lng, rad)
	query := r.URL.Query()
	lat, _ := strconv.ParseFloat(query.Get("lat"), 64)
	lng, _ := strconv.ParseFloat(query.Get("lng"), 64)
	rad, _ := strconv.ParseFloat(query.Get("rad"), 64)

	// Default radius if missing
	if rad == 0 {
		rad = 10.0 // 10km
	}

	// 3. Call the Repository (The Logic)
	shifts, err := h.Repo.FindNearby(r.Context(), lat, lng, rad)
	if err != nil {
		http.Error(w, "Failed to search shifts", http.StatusInternalServerError)
		return
	}

	// 4. Send Response
	json.NewEncoder(w).Encode(shifts)
}