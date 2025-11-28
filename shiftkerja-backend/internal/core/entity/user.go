package entity

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "json:-" means never send this back to the frontend!
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"` // 'worker' or 'business'
	CreatedAt    time.Time `json:"created_at"`
}