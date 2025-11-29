package dto

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=worker business admin"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
	Email string `json:"email"`
}
