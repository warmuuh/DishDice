package services

import (
	"fmt"
	"time"

	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	ticketRepo *repository.TicketRepository
	jwtSecret  string
}

func NewAuthService(userRepo *repository.UserRepository, ticketRepo *repository.TicketRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		ticketRepo: ticketRepo,
		jwtSecret:  jwtSecret,
	}
}

func (s *AuthService) Register(email, password, language, ticket string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Validate ticket if provided
	autoApprove := false
	if ticket != "" {
		ticketData, err := s.ticketRepo.GetByToken(ticket)
		if err != nil {
			return nil, fmt.Errorf("failed to validate ticket: %w", err)
		}
		if ticketData == nil {
			return nil, fmt.Errorf("invalid registration ticket")
		}
		if ticketData.IsUsed {
			return nil, fmt.Errorf("registration ticket has already been used")
		}
		if time.Now().After(ticketData.ExpiresAt) {
			return nil, fmt.Errorf("registration ticket has expired")
		}
		autoApprove = true
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with language preference and ticket auto-approval
	user, err := s.userRepo.CreateWithLanguageAndTicket(email, string(hashedPassword), language, autoApprove)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Mark ticket as used if applicable
	if ticket != "" && autoApprove {
		err = s.ticketRepo.MarkTicketUsed(ticket, user.ID)
		if err != nil {
			// Log warning but don't fail registration since user was already created
			fmt.Printf("Warning: failed to mark ticket as used: %v\n", err)
		}
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	// Check user status
	if user.Status == models.StatusPending {
		return nil, "", fmt.Errorf("account is pending approval")
	}
	if user.Status == models.StatusRejected {
		return nil, "", fmt.Errorf("account has been rejected")
	}
	if user.Status != models.StatusApproved {
		return nil, "", fmt.Errorf("account is not active")
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"status":  user.Status,
		"exp":     time.Now().Add(24 * 7 * time.Hour).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

type TokenClaims struct {
	UserID string
	Role   string
	Status string
}

func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	role, _ := claims["role"].(string)
	status, _ := claims["status"].(string)

	return &TokenClaims{
		UserID: userID,
		Role:   role,
		Status: status,
	}, nil
}
