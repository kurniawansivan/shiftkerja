package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"shiftkerja-backend/internal/core/entity"

	"github.com/redis/go-redis/v9"
)

type RedisGeoRepository struct {
	Client *redis.Client
}

func NewRedisGeoRepo(client *redis.Client) *RedisGeoRepository {
	return &RedisGeoRepository{Client: client}
}

// AddShift saves the Shift metadata AND its location
func (r *RedisGeoRepository) AddShift(ctx context.Context, shift entity.Shift) error {
	// 1. Store the details as JSON (so we can read the Title/Pay later)
	data, _ := json.Marshal(shift)
	
	// We use a simple Key-Value pair for details: "shift:123" -> JSON
	key := fmt.Sprintf("shift:%d", shift.ID) 
	if err := r.Client.Set(ctx, key, data, 0).Err(); err != nil {
		return err
	}

	// 2. Store the Location in the Geospatial Index
	// "shifts_geo" is the key for our map index
	cmd := r.Client.GeoAdd(ctx, "shifts_geo", &redis.GeoLocation{
		Name:      fmt.Sprintf("%d", shift.ID),
		Longitude: shift.Lng,
		Latitude:  shift.Lat,
	})
	
	return cmd.Err()
}

// FindNearby returns all shifts within 'km' radius
func (r *RedisGeoRepository) FindNearby(ctx context.Context, lat, lng, km float64) ([]entity.Shift, error) {
	// 1. Ask Redis: "Give me IDs within X km"
	locations, err := r.Client.GeoSearch(ctx, "shifts_geo", &redis.GeoSearchQuery{
		Longitude:  lng,
		Latitude:   lat,
		Radius:     km,
		RadiusUnit: "km",
	}).Result()

	if err != nil {
		return nil, err
	}

	// 2. Resolve those IDs into full Shift objects
	var shifts []entity.Shift
	for _, loc := range locations {
		data, err := r.Client.Get(ctx, fmt.Sprintf("shift:%s", loc)).Bytes()
		if err == nil {
			var s entity.Shift
			if json.Unmarshal(data, &s) == nil {
				shifts = append(shifts, s)
			}
		}
	}

	return shifts, nil
}

// RemoveShift removes a shift from the geo index
func (r *RedisGeoRepository) RemoveShift(ctx context.Context, shiftID int64) error {
	// 1. Remove from geo index
	key := fmt.Sprintf("%d", shiftID)
	if err := r.Client.ZRem(ctx, "shifts_geo", key).Err(); err != nil {
		return err
	}
	
	// 2. Remove the metadata
	dataKey := fmt.Sprintf("shift:%d", shiftID)
	if err := r.Client.Del(ctx, dataKey).Err(); err != nil {
		return err
	}
	
	return nil
}