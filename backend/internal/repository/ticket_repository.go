package repository

import (
	"database/sql"
	"time"

	"github.com/dishdice/backend/internal/models"
	"github.com/google/uuid"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) CreateTicket(createdBy string, expiresAt time.Time) (*models.RegistrationTicket, error) {
	token := uuid.New().String()

	query := `
		INSERT INTO registration_tickets (token, created_by, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, token, created_by, used_by, created_at, expires_at, used_at, is_used
	`

	var ticket models.RegistrationTicket
	err := r.db.QueryRow(query, token, createdBy, expiresAt).Scan(
		&ticket.ID,
		&ticket.Token,
		&ticket.CreatedBy,
		&ticket.UsedBy,
		&ticket.CreatedAt,
		&ticket.ExpiresAt,
		&ticket.UsedAt,
		&ticket.IsUsed,
	)

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketRepository) GetByToken(token string) (*models.RegistrationTicket, error) {
	query := `
		SELECT id, token, created_by, used_by, created_at, expires_at, used_at, is_used
		FROM registration_tickets
		WHERE token = $1
	`

	var ticket models.RegistrationTicket
	err := r.db.QueryRow(query, token).Scan(
		&ticket.ID,
		&ticket.Token,
		&ticket.CreatedBy,
		&ticket.UsedBy,
		&ticket.CreatedAt,
		&ticket.ExpiresAt,
		&ticket.UsedAt,
		&ticket.IsUsed,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketRepository) MarkTicketUsed(token, usedBy string) error {
	query := `
		UPDATE registration_tickets
		SET is_used = true, used_by = $1, used_at = NOW()
		WHERE token = $2 AND is_used = false
	`

	result, err := r.db.Exec(query, usedBy, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TicketRepository) DeleteExpiredTickets() (int64, error) {
	query := `
		DELETE FROM registration_tickets
		WHERE expires_at < NOW() AND is_used = false
	`

	result, err := r.db.Exec(query)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
