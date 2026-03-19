package models

import "time"

type RegistrationTicket struct {
	ID        string     `json:"id"`
	Token     string     `json:"token"`
	CreatedBy string     `json:"created_by"`
	UsedBy    *string    `json:"used_by,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	IsUsed    bool       `json:"is_used"`
}

type CreateTicketResponse struct {
	Ticket           RegistrationTicket `json:"ticket"`
	RegistrationLink string             `json:"registration_link"`
}

type ValidateTicketResponse struct {
	Valid     bool   `json:"valid"`
	Message   string `json:"message,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}
