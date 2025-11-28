package service

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Secret Key (In production, load this from .env!)
var jwtSecret = []byte("SUPER_SECRET_KEY_DO_NOT_SHARE")

func GenerateToken(userID int64, role string) (string, error) {
	// Define the "Claims" (Data inside the token)
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	}

	// Create and Sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}