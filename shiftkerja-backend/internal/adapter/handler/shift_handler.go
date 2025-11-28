package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"shiftkerja-backend/internal/core/entity"
	"shiftkerja-backend/internal/core/service"
)

type ShiftHandler struct {
	Service *service.ShiftService
}

// Constructor using service layer (Clean Architecture)
func NewShiftHandler(svc *service.ShiftService) *ShiftHandler {
	return &ShiftHandler{
		Service: svc,
	}
}

// GetNearby returns shifts within radius (for map)
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

	// Call service layer
	shifts, err := h.Service.GetNearbyShifts(r.Context(), lat, lng, rad)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(shifts)
}

// Create handles shift creation (Business only)
func (h *ShiftHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Security Check
	role := r.Context().Value("role").(string)
	userID := r.Context().Value("user_id").(float64)

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

	// Set the owner from the token
	req.OwnerID = int64(userID)

	// 3. Call service layer (handles dual-write)
	if err := h.Service.CreateShift(r.Context(), &req); err != nil {
		fmt.Printf("❌ Create Shift Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// ApplyRequest defines the request body for applying to a shift
type ApplyRequest struct {
	ShiftID int64 `json:"shift_id"`
}

// Apply handles worker application to a shift
func (h *ShiftHandler) Apply(w http.ResponseWriter, r *http.Request) {
	// 1. Authentication check
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "worker" {
		http.Error(w, "Only workers can apply", http.StatusForbidden)
		return
	}

	// 2. Parse request
	var req ApplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 3. Call service layer
	err := h.Service.ApplyForShift(r.Context(), req.ShiftID, userID)
	if err != nil {
		fmt.Printf("❌ Apply Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Applied successfully"})
}

// GetMyShifts returns all shifts posted by the business owner
func (h *ShiftHandler) GetMyShifts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "business" && role != "admin" {
		http.Error(w, "Only businesses can view their shifts", http.StatusForbidden)
		return
	}

	shifts, err := h.Service.GetMyShifts(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve shifts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(shifts)
}

// GetMyApplications returns all applications for a worker
func (h *ShiftHandler) GetMyApplications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "worker" {
		http.Error(w, "Only workers can view their applications", http.StatusForbidden)
		return
	}

	applications, err := h.Service.GetMyApplications(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(applications)
}

// GetShiftApplications returns all applications for a specific shift (business owner only)
func (h *ShiftHandler) GetShiftApplications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	userID := int64(r.Context().Value("user_id").(float64))
	
	// Parse shift ID from query param
	shiftID, err := strconv.ParseInt(r.URL.Query().Get("shift_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid shift_id", http.StatusBadRequest)
		return
	}

	applications, err := h.Service.GetShiftApplications(r.Context(), shiftID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(applications)
}

// UpdateApplicationStatusRequest defines the request body for updating application status
type UpdateApplicationStatusRequest struct {
	ApplicationID int64  `json:"application_id"`
	Status        string `json:"status"` // ACCEPTED or REJECTED
}

// UpdateApplicationStatus handles accepting/rejecting applications
func (h *ShiftHandler) UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "business" && role != "admin" {
		http.Error(w, "Only businesses can update application status", http.StatusForbidden)
		return
	}

	var req UpdateApplicationStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateApplicationStatus(r.Context(), req.ApplicationID, userID, req.Status)
	if err != nil {
		fmt.Printf("❌ Update Status Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Updated successfully"})
}