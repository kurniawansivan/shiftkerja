package port

import (
	"context"
	"shiftkerja-backend/internal/core/entity"
)

// GeoRepository defines the contract for geospatial operations
type GeoRepository interface {
	AddShift(ctx context.Context, shift entity.Shift) error
	FindNearby(ctx context.Context, lat, lng, radiusKm float64) ([]entity.Shift, error)
	RemoveShift(ctx context.Context, shiftID int64) error
}
