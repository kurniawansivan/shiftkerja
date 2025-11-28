package handler

import (
	"context"
	"net/http"
	"strings"
	"shiftkerja-backend/internal/core/service"
)

// AuthMiddleware wraps a standard http.HandlerFunc with security checks
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the Auth Header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		// 2. Remove "Bearer " prefix
		// Format should be: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid Token Format (Missing 'Bearer')", http.StatusUnauthorized)
			return
		}

		// 3. Validate Token
		claims, err := service.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or Expired Token", http.StatusUnauthorized)
			return
		}

		// 4. (Optional) Pass user info down to the handler
		// We store the user_id in the request Context so the next handler knows WHO called it
		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		ctx = context.WithValue(ctx, "role", claims["role"])

		// 5. Pass to the next handler (The "Real" logic)
		next(w, r.WithContext(ctx))
	}
}