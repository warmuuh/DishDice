package services

import (
	"fmt"
	"time"

	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/repository"
)

type TicketService struct {
	ticketRepo  *repository.TicketRepository
	frontendURL string
}

func NewTicketService(ticketRepo *repository.TicketRepository, frontendURL string) *TicketService {
	return &TicketService{
		ticketRepo:  ticketRepo,
		frontendURL: frontendURL,
	}
}

func (s *TicketService) CreateTicket(adminID string) (*models.CreateTicketResponse, error) {
	expiresAt := time.Now().Add(14 * 24 * time.Hour) // 2 weeks

	ticket, err := s.ticketRepo.CreateTicket(adminID, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	registrationLink := fmt.Sprintf("%s/register?ticket=%s", s.frontendURL, ticket.Token)

	return &models.CreateTicketResponse{
		Ticket:           *ticket,
		RegistrationLink: registrationLink,
	}, nil
}

func (s *TicketService) ValidateTicket(token string) (*models.ValidateTicketResponse, error) {
	ticket, err := s.ticketRepo.GetByToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	if ticket == nil {
		return &models.ValidateTicketResponse{
			Valid:   false,
			Message: "Ticket not found",
		}, nil
	}

	if ticket.IsUsed {
		return &models.ValidateTicketResponse{
			Valid:   false,
			Message: "Ticket has already been used",
		}, nil
	}

	if time.Now().After(ticket.ExpiresAt) {
		return &models.ValidateTicketResponse{
			Valid:   false,
			Message: "Ticket has expired",
		}, nil
	}

	return &models.ValidateTicketResponse{
		Valid:     true,
		ExpiresAt: ticket.ExpiresAt.Format(time.RFC3339),
	}, nil
}
