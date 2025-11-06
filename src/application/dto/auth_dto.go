package dto

import "time"

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse represents the user information in API responses
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Department string   `json:"department"`
	CreatedAt time.Time `json:"created_at"`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents the refresh token response payload
type RefreshTokenResponse struct {
	Token string `json:"token"`
}