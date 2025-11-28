package entity

import "time"

type Application struct {
	ID        int64     `json:"id"`
	ShiftID   int64     `json:"shift_id"`
	WorkerID  int64     `json:"worker_id"`
	Status    string    `json:"status"` // PENDING, ACCEPTED, REJECTED
	CreatedAt time.Time `json:"created_at"`
	
	// Populated via JOIN queries
	ShiftTitle   string  `json:"shift_title,omitempty"`
	ShiftPayRate float64 `json:"shift_pay_rate,omitempty"`
	WorkerName   string  `json:"worker_name,omitempty"`
	WorkerEmail  string  `json:"worker_email,omitempty"`
}
