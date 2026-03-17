package services

import (
	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/repository"
)

type AdminService struct {
	userRepo *repository.UserRepository
}

func NewAdminService(userRepo *repository.UserRepository) *AdminService {
	return &AdminService{userRepo: userRepo}
}

func (s *AdminService) GetAllUsers() ([]*models.AdminUserListItem, error) {
	return s.userRepo.GetAllUsersForAdmin()
}

func (s *AdminService) GetPendingUsers() ([]*models.AdminUserListItem, error) {
	return s.userRepo.GetUsersByStatus(models.StatusPending)
}

func (s *AdminService) ApproveUser(userID string) error {
	return s.userRepo.UpdateUserStatus(userID, models.StatusApproved)
}

func (s *AdminService) RejectUser(userID string) error {
	return s.userRepo.UpdateUserStatus(userID, models.StatusRejected)
}
