package repository

import (
	"database/sql"
	"fmt"

	"github.com/dishdice/backend/internal/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, passwordHash string) (*models.User, error) {
	return r.CreateWithLanguage(email, passwordHash, "en")
}

func (r *UserRepository) CreateWithLanguage(email, passwordHash, language string) (*models.User, error) {
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		Language:     language,
	}

	// Check if this is the first user
	var userCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// First user becomes admin and is auto-approved
	role := models.RoleUser
	status := models.StatusPending
	if userCount == 0 {
		role = models.RoleAdmin
		status = models.StatusApproved
	}

	query := `
		INSERT INTO users (id, email, password_hash, language, role, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err = r.db.QueryRow(query, user.ID, user.Email, user.PasswordHash, user.Language, role, status).
		Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.Role = role
	user.Status = status

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, preferences, language, role, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Preferences,
		&user.Language,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, preferences, language, role, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Preferences,
		&user.Language,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *UserRepository) UpdatePreferences(id, preferences, language string) error {
	query := `
		UPDATE users
		SET preferences = $1, language = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.db.Exec(query, preferences, language, id)
	if err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) GetUsersByStatus(status string) ([]*models.AdminUserListItem, error) {
	query := `
		SELECT id, email, role, status, language, created_at
		FROM users
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by status: %w", err)
	}
	defer rows.Close()

	var users []*models.AdminUserListItem
	for rows.Next() {
		user := &models.AdminUserListItem{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Role,
			&user.Status,
			&user.Language,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetAllUsersForAdmin() ([]*models.AdminUserListItem, error) {
	query := `
		SELECT id, email, role, status, language, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []*models.AdminUserListItem
	for rows.Next() {
		user := &models.AdminUserListItem{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Role,
			&user.Status,
			&user.Language,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) UpdateUserStatus(userID, status string) error {
	query := `
		UPDATE users
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, status, userID)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
