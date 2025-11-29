package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"shiftkerja-backend/internal/core/dto"
	"shiftkerja-backend/internal/core/entity"
	"shiftkerja-backend/internal/core/service"
	"shiftkerja-backend/pkg/util"
)

type ShiftHandler struct {
	Service *service.ShiftService
	Hub     *Hub
}

// Constructor using service layer (Clean Architecture)
func NewShiftHandler(svc *service.ShiftService, hub *Hub) *ShiftHandler {
	return &ShiftHandler{
		Service: svc,
		Hub:     hub,
	}
}

// GetNearby returns shifts within radius (for map)
func (h *ShiftHandler) GetNearby(w http.ResponseWriter, r *http.Request) {
	// Parse Query Params
	q := r.URL.Query()
	lat, err := strconv.ParseFloat(q.Get("lat"), 64)
	if err != nil || lat < -90 || lat > 90 {
		util.RespondBadRequest(w, "Invalid latitude: must be between -90 and 90")
		return
	}
	
	lng, err := strconv.ParseFloat(q.Get("lng"), 64)
	if err != nil || lng < -180 || lng > 180 {
		util.RespondBadRequest(w, "Invalid longitude: must be between -180 and 180")
		return
	}
	
	rad, _ := strconv.ParseFloat(q.Get("rad"), 64)
	if rad <= 0 {
		rad = 10 // Default 10km radius
	}
	if rad > 100 {
		util.RespondBadRequest(w, "Radius cannot exceed 100km")
		return
	}

	// Call service layer
	shifts, err := h.Service.GetNearbyShifts(r.Context(), lat, lng, rad)
	if err != nil {
		fmt.Printf("❌ GetNearby Error: %v\n", err)
		util.RespondInternalError(w, "Failed to search for shifts")
		return
	}
	
	// Return empty array if no shifts found
	if shifts == nil {
		shifts = []entity.Shift{}
	}
	
	util.RespondJSON(w, http.StatusOK, shifts)
}

// Create handles shift creation (Business only)
func (h *ShiftHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Security Check
	role := r.Context().Value("role").(string)
	userID := r.Context().Value("user_id").(float64)

	if role != "business" && role != "admin" {
		util.RespondForbidden(w, "Only businesses can post shifts")
		return
	}

	// 2. Parse Input
	var req dto.CreateShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondBadRequest(w, "Invalid JSON format")
		return
	}

	// 3. Validate input
	if req.Title == "" {
		util.RespondBadRequest(w, "Title is required")
		return
	}
	if req.PayRate <= 0 {
		util.RespondBadRequest(w, "Pay rate must be greater than 0")
		return
	}
	if req.Lat < -90 || req.Lat > 90 {
		util.RespondBadRequest(w, "Latitude must be between -90 and 90")
		return
	}
	if req.Lng < -180 || req.Lng > 180 {
		util.RespondBadRequest(w, "Longitude must be between -180 and 180")
		return
	}

	// 4. Convert to entity
	shift := &entity.Shift{
		OwnerID:     int64(userID),
		Title:       req.Title,
		Description: req.Description,
		PayRate:     req.PayRate,
		Lat:         req.Lat,
		Lng:         req.Lng,
		Status:      "OPEN",
	}

	// 5. Call service layer (handles dual-write)
	if err := h.Service.CreateShift(r.Context(), shift); err != nil {
		fmt.Printf("❌ Create Shift Error: %v\n", err)
		util.RespondInternalError(w, err.Error())
		return
	}

	// 6. BROADCAST TO WEBSOCKET (Live Updates)
	if h.Hub != nil {
		broadcastMsg := map[string]interface{}{
			"type":     "shift_created",
			"id":       shift.ID,
			"title":    shift.Title,
			"lat":      shift.Lat,
			"lng":      shift.Lng,
			"pay_rate": shift.PayRate,
			"status":   shift.Status,
		}
		h.Hub.Broadcast(broadcastMsg)
		fmt.Printf("📡 Broadcasted shift creation: %s\n", shift.Title)
	}

	util.RespondCreated(w, "Shift created successfully", shift)
}

// Apply handles worker application to a shift
func (h *ShiftHandler) Apply(w http.ResponseWriter, r *http.Request) {
	// 1. Authentication check
	userIDFloat, ok := r.Context().Value("user_id").(float64)
	if !ok {
		util.RespondUnauthorized(w, "Invalid authentication token")
		return
	}
	userID := int64(userIDFloat)
	
	role, ok := r.Context().Value("role").(string)
	if !ok {
		util.RespondUnauthorized(w, "Invalid role in token")
		return
	}

	if role != "worker" {
		util.RespondForbidden(w, "Only workers can apply to shifts")
		return
	}

	// 2. Parse request
	var req dto.ApplyShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondBadRequest(w, "Invalid JSON format: "+err.Error())
		return
	}

	// 3. Validate input
	if req.ShiftID <= 0 {
		util.RespondBadRequest(w, "Invalid shift_id: must be greater than 0")
		return
	}

	fmt.Printf("🔄 Worker %d applying for shift %d\n", userID, req.ShiftID)

	// 4. Call service layer
	err := h.Service.ApplyForShift(r.Context(), req.ShiftID, userID)
	if err != nil {
		fmt.Printf("❌ Apply Error: %v\n", err)
		
		// Send appropriate error response based on error type
		switch err {
		case service.ErrShiftNotFound:
			util.RespondNotFound(w, "Shift not found")
		case service.ErrApplicationExists:
			util.RespondBadRequest(w, "You have already applied to this shift")
		default:
			util.RespondBadRequest(w, err.Error())
		}
		return
	}

	fmt.Printf("✅ Application successful: Worker %d -> Shift %d\n", userID, req.ShiftID)

	// 5. BROADCAST DETAILED APPLICATION EVENT
	if h.Hub != nil {
		// Get latest application details including worker info
		applications, _ := h.Service.GetMyApplications(r.Context(), userID)
		
		broadcastMsg := map[string]interface{}{
			"type":      "new_application",
			"shift_id":  req.ShiftID,
			"worker_id": userID,
			"status":    "PENDING",
		}
		
		// Find the application we just created and add details
		for _, app := range applications {
			if app.ShiftID == req.ShiftID {
				broadcastMsg["application_id"] = app.ID
				broadcastMsg["shift_title"] = app.ShiftTitle
				broadcastMsg["shift_pay_rate"] = app.ShiftPayRate
				broadcastMsg["created_at"] = app.CreatedAt
				break
			}
		}
		
		h.Hub.Broadcast(broadcastMsg)
		fmt.Printf("📡 Broadcasted new application: Worker %d -> Shift %d\n", userID, req.ShiftID)
	}

	util.RespondSuccess(w, "Application submitted successfully", map[string]interface{}{
		"shift_id":  req.ShiftID,
		"worker_id": userID,
		"status":    "PENDING",
	})
}

// GetMyShifts returns all shifts posted by the business owner
func (h *ShiftHandler) GetMyShifts(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "business" && role != "admin" {
		util.RespondForbidden(w, "Only businesses can view their shifts")
		return
	}

	shifts, err := h.Service.GetMyShifts(r.Context(), userID)
	if err != nil {
		fmt.Printf("❌ GetMyShifts Error: %v\n", err)
		util.RespondInternalError(w, "Failed to retrieve shifts")
		return
	}

	// Return empty array if no shifts
	if shifts == nil {
		shifts = []entity.Shift{}
	}

	util.RespondJSON(w, http.StatusOK, shifts)
}

// GetMyApplications returns all applications for a worker
func (h *ShiftHandler) GetMyApplications(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "worker" {
		util.RespondForbidden(w, "Only workers can view their applications")
		return
	}

	applications, err := h.Service.GetMyApplications(r.Context(), userID)
	if err != nil {
		fmt.Printf("❌ GetMyApplications Error: %v\n", err)
		util.RespondInternalError(w, "Failed to retrieve applications")
		return
	}

	// Return empty array if no applications
	if applications == nil {
		applications = []entity.Application{}
	}

	util.RespondJSON(w, http.StatusOK, applications)
}

// GetShiftApplications returns all applications for a specific shift (business owner only)
func (h *ShiftHandler) GetShiftApplications(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(float64))
	
	// Parse shift ID from query param
	shiftID, err := strconv.ParseInt(r.URL.Query().Get("shift_id"), 10, 64)
	if err != nil || shiftID <= 0 {
		util.RespondBadRequest(w, "Invalid shift_id: must be a positive integer")
		return
	}

	applications, err := h.Service.GetShiftApplications(r.Context(), shiftID, userID)
	if err != nil {
		fmt.Printf("❌ GetShiftApplications Error: %v\n", err)
		if err == service.ErrUnauthorized {
			util.RespondForbidden(w, "You don't have permission to view these applications")
			return
		}
		if err == service.ErrShiftNotFound {
			util.RespondNotFound(w, "Shift not found")
			return
		}
		util.RespondInternalError(w, "Failed to retrieve applications")
		return
	}

	// Return empty array if no applications
	if applications == nil {
		applications = []entity.Application{}
	}

	util.RespondJSON(w, http.StatusOK, applications)
}

// UpdateApplicationStatus handles accepting/rejecting applications
func (h *ShiftHandler) UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "business" && role != "admin" {
		util.RespondForbidden(w, "Only businesses can update application status")
		return
	}

	var req dto.UpdateApplicationStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondBadRequest(w, "Invalid JSON format")
		return
	}

	// Validate input
	if req.ApplicationID <= 0 {
		util.RespondBadRequest(w, "Invalid application_id")
		return
	}
	if req.Status != "ACCEPTED" && req.Status != "REJECTED" {
		util.RespondBadRequest(w, "Status must be either ACCEPTED or REJECTED")
		return
	}

	err := h.Service.UpdateApplicationStatus(r.Context(), req.ApplicationID, userID, req.Status)
	if err != nil {
		fmt.Printf("❌ Update Status Error: %v\n", err)
		
		switch err {
		case service.ErrUnauthorized:
			util.RespondForbidden(w, "You don't have permission to update this application")
		case service.ErrShiftNotFound:
			util.RespondNotFound(w, "Shift not found")
		case service.ErrInvalidStatus:
			util.RespondBadRequest(w, "Invalid status transition")
		default:
			util.RespondBadRequest(w, err.Error())
		}
		return
	}

	// BROADCAST APPLICATION STATUS UPDATE
	if h.Hub != nil {
		broadcastMsg := map[string]interface{}{
			"type":           "application_status_updated",
			"application_id": req.ApplicationID,
			"new_status":     req.Status,
			"updated_by":     userID,
		}
		h.Hub.Broadcast(broadcastMsg)
		fmt.Printf("📡 Broadcasted status update: Application %d -> %s\n", req.ApplicationID, req.Status)
	}

	util.RespondSuccess(w, "Application status updated successfully", map[string]interface{}{
		"application_id": req.ApplicationID,
		"new_status":     req.Status,
	})
}

// UpdateShift updates an existing shift
func (h *ShiftHandler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "business" {
		util.RespondForbidden(w, "Only businesses can update shifts")
		return
	}

	var req dto.UpdateShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondBadRequest(w, "Invalid JSON format")
		return
	}

	// Validate input
	if req.ID <= 0 {
		util.RespondBadRequest(w, "Invalid shift ID")
		return
	}
	if req.Title == "" {
		util.RespondBadRequest(w, "Title is required")
		return
	}
	if req.PayRate <= 0 {
		util.RespondBadRequest(w, "Pay rate must be greater than 0")
		return
	}
	if req.Status != "OPEN" && req.Status != "FILLED" {
		util.RespondBadRequest(w, "Status must be either OPEN or FILLED")
		return
	}

	// Convert to entity
	shift := &entity.Shift{
		ID:          req.ID,
		OwnerID:     userID,
		Title:       req.Title,
		Description: req.Description,
		PayRate:     req.PayRate,
		Lat:         req.Lat,
		Lng:         req.Lng,
		Status:      req.Status,
	}

	err := h.Service.UpdateShift(r.Context(), shift, userID)
	if err != nil {
		fmt.Printf("❌ Update Shift Error: %v\n", err)
		
		switch err {
		case service.ErrUnauthorized:
			util.RespondForbidden(w, "You don't have permission to update this shift")
		case service.ErrShiftNotFound:
			util.RespondNotFound(w, "Shift not found")
		default:
			util.RespondBadRequest(w, err.Error())
		}
		return
	}

	util.RespondSuccess(w, "Shift updated successfully", shift)
}

// DeleteShift deletes a shift
func (h *ShiftHandler) DeleteShift(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "business" {
		util.RespondForbidden(w, "Only businesses can delete shifts")
		return
	}

	shiftIDStr := r.URL.Query().Get("shift_id")
	if shiftIDStr == "" {
		util.RespondBadRequest(w, "shift_id is required")
		return
	}

	shiftID, err := strconv.ParseInt(shiftIDStr, 10, 64)
	if err != nil || shiftID <= 0 {
		util.RespondBadRequest(w, "Invalid shift_id: must be a positive integer")
		return
	}

	err = h.Service.DeleteShift(r.Context(), shiftID, userID)
	if err != nil {
		fmt.Printf("❌ Delete Shift Error: %v\n", err)
		
		switch err {
		case service.ErrUnauthorized:
			util.RespondForbidden(w, "You don't have permission to delete this shift")
		case service.ErrShiftNotFound:
			util.RespondNotFound(w, "Shift not found")
		default:
			util.RespondInternalError(w, "Failed to delete shift")
		}
		return
	}

	fmt.Printf("✅ Shift %d deleted successfully by user %d\n", shiftID, userID)
	util.RespondSuccess(w, "Shift deleted successfully", map[string]interface{}{
		"shift_id": shiftID,
	})
}

// DeleteApplication allows a worker to withdraw their application
func (h *ShiftHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(float64))
	role := r.Context().Value("role").(string)

	if role != "worker" {
		util.RespondForbidden(w, "Only workers can delete their applications")
		return
	}

	appIDStr := r.URL.Query().Get("application_id")
	if appIDStr == "" {
		util.RespondBadRequest(w, "application_id is required")
		return
	}

	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil || appID <= 0 {
		util.RespondBadRequest(w, "Invalid application_id: must be a positive integer")
		return
	}

	err = h.Service.DeleteApplication(r.Context(), appID, userID)
	if err != nil {
		fmt.Printf("❌ Delete Application Error: %v\n", err)
		
		switch err {
		case service.ErrUnauthorized:
			util.RespondForbidden(w, "You don't have permission to delete this application")
		default:
			util.RespondBadRequest(w, err.Error())
		}
		return
	}

	fmt.Printf("✅ Application %d withdrawn by worker %d\n", appID, userID)
	util.RespondSuccess(w, "Application withdrawn successfully", map[string]interface{}{
		"application_id": appID,
	})
}