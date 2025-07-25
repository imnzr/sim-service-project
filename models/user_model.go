package models

import "time"

// User merepresentasikan pengguna di sistem
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Register user payload untuk pendaftaran pengguna baru
type RegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login request payload untuk login pengguna
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Token response setelah login/refresh token sukses
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// User profile resposne untuk detail profile pengguna
type UserProfileResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
