package model

import (
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	Activated      bool      `json:"activated"`
	ActivationLink string    `json:"activation_link"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
