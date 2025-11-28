package entity

import "time"

type Shift struct {
	ID          int64     `json:"id"`       
	OwnerID     int64     `json:"owner_id"` 
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PayRate     float64   `json:"pay_rate"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
	Status      string    `json:"status"`    
	CreatedAt   time.Time `json:"created_at"`
}