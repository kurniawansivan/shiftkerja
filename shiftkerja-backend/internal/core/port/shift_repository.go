package port

import (
	"context"
	"shiftkerja-backend/internal/core/entity"
)

// ShiftRepository defines the contract for shift data access
type ShiftRepository interface {
	CreateShift(ctx context.Context, shift *entity.Shift) error
	GetShiftByID(ctx context.Context, id int64) (*entity.Shift, error)
	GetShiftsByOwner(ctx context.Context, ownerID int64) ([]entity.Shift, error)
	UpdateShiftStatus(ctx context.Context, id int64, status string) error
	
	// Application methods
	ApplyForShift(ctx context.Context, shiftID, workerID int64) error
	GetApplicationsByWorker(ctx context.Context, workerID int64) ([]entity.Application, error)
	GetApplicationsByShift(ctx context.Context, shiftID int64) ([]entity.Application, error)
	UpdateApplicationStatus(ctx context.Context, applicationID int64, status string) error
	GetApplicationByID(ctx context.Context, id int64) (*entity.Application, error)
}
