package service

import (
	"context"
	"errors"
	"fmt"
	"shiftkerja-backend/internal/core/entity"
	"shiftkerja-backend/internal/core/port"
)

var (
	ErrUnauthorized       = errors.New("unauthorized action")
	ErrShiftNotFound      = errors.New("shift not found")
	ErrApplicationExists  = errors.New("already applied to this shift")
	ErrInvalidStatus      = errors.New("invalid status transition")
)

type ShiftService struct {
	shiftRepo port.ShiftRepository
	geoRepo   port.GeoRepository
}

func NewShiftService(shiftRepo port.ShiftRepository, geoRepo port.GeoRepository) *ShiftService {
	return &ShiftService{
		shiftRepo: shiftRepo,
		geoRepo:   geoRepo,
	}
}

// CreateShift handles shift creation with dual-write to Postgres and Redis
func (s *ShiftService) CreateShift(ctx context.Context, shift *entity.Shift) error {
	// 1. Validate business rules
	if shift.PayRate <= 0 {
		return errors.New("pay rate must be positive")
	}
	if shift.Title == "" {
		return errors.New("title is required")
	}
	
	// 2. Save to Postgres (source of truth)
	if err := s.shiftRepo.CreateShift(ctx, shift); err != nil {
		return fmt.Errorf("failed to create shift: %w", err)
	}
	
	// 3. Sync to Redis (geo index)
	if err := s.geoRepo.AddShift(ctx, *shift); err != nil {
		// Log but don't fail - data is in Postgres
		fmt.Printf("⚠️ Redis sync warning: %v\n", err)
	}
	
	return nil
}

// GetNearbyShifts retrieves shifts within radius
func (s *ShiftService) GetNearbyShifts(ctx context.Context, lat, lng, radiusKm float64) ([]entity.Shift, error) {
	return s.geoRepo.FindNearby(ctx, lat, lng, radiusKm)
}

// ApplyForShift handles worker application with validation
func (s *ShiftService) ApplyForShift(ctx context.Context, shiftID, workerID int64) error {
	// 1. Check if shift exists
	shift, err := s.shiftRepo.GetShiftByID(ctx, shiftID)
	if err != nil {
		return ErrShiftNotFound
	}
	
	// 2. Check if shift is still open
	if shift.Status != "OPEN" {
		return errors.New("shift is no longer available")
	}
	
	// 3. Apply
	if err := s.shiftRepo.ApplyForShift(ctx, shiftID, workerID); err != nil {
		return fmt.Errorf("failed to apply: %w", err)
	}
	
	return nil
}

// GetMyShifts retrieves shifts posted by a business owner
func (s *ShiftService) GetMyShifts(ctx context.Context, ownerID int64) ([]entity.Shift, error) {
	return s.shiftRepo.GetShiftsByOwner(ctx, ownerID)
}

// GetMyApplications retrieves applications for a worker
func (s *ShiftService) GetMyApplications(ctx context.Context, workerID int64) ([]entity.Application, error) {
	return s.shiftRepo.GetApplicationsByWorker(ctx, workerID)
}

// GetShiftApplications retrieves all applications for a shift (business owner only)
func (s *ShiftService) GetShiftApplications(ctx context.Context, shiftID, requesterID int64) ([]entity.Application, error) {
	// Verify ownership
	shift, err := s.shiftRepo.GetShiftByID(ctx, shiftID)
	if err != nil {
		return nil, ErrShiftNotFound
	}
	
	if shift.OwnerID != requesterID {
		return nil, ErrUnauthorized
	}
	
	return s.shiftRepo.GetApplicationsByShift(ctx, shiftID)
}

// UpdateApplicationStatus handles accepting/rejecting applications
func (s *ShiftService) UpdateApplicationStatus(ctx context.Context, applicationID, businessID int64, newStatus string) error {
	// 1. Validate status
	if newStatus != "ACCEPTED" && newStatus != "REJECTED" {
		return ErrInvalidStatus
	}
	
	// 2. Get application details
	app, err := s.shiftRepo.GetApplicationByID(ctx, applicationID)
	if err != nil {
		return errors.New("application not found")
	}
	
	// 3. Verify the requester owns the shift
	shift, err := s.shiftRepo.GetShiftByID(ctx, app.ShiftID)
	if err != nil {
		return ErrShiftNotFound
	}
	
	if shift.OwnerID != businessID {
		return ErrUnauthorized
	}
	
	// 4. Update application status
	if err := s.shiftRepo.UpdateApplicationStatus(ctx, applicationID, newStatus); err != nil {
		return err
	}
	
	// 5. If accepted, update shift status to FILLED
	if newStatus == "ACCEPTED" {
		if err := s.shiftRepo.UpdateShiftStatus(ctx, app.ShiftID, "FILLED"); err != nil {
			return err
		}
		
		// Remove from Redis geo index
		if err := s.geoRepo.RemoveShift(ctx, app.ShiftID); err != nil {
			fmt.Printf("⚠️ Redis remove warning: %v\n", err)
		}
	}
	
	return nil
}
