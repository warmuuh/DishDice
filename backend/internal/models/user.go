package models

import "time"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"

	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Preferences  *string   `json:"preferences"`
	Language     string    `json:"language"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdatePreferencesRequest struct {
	Preferences string `json:"preferences"`
	Language    string `json:"language"`
}

type RegisterResponse struct {
	Message string  `json:"message,omitempty"`
	Status  string  `json:"status"`
	Token   *string `json:"token,omitempty"`
	User    User    `json:"user"`
}

type AdminUserListItem struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
}
