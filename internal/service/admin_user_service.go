package service

import (
	"errors"
	"go-ticket/internal/domain"
)

type adminUserService struct {
	userRepo domain.UserRepository
}

func NewAdminUserService(userRepo domain.UserRepository) domain.AdminUserService {
	return &adminUserService{
		userRepo: userRepo,
	}
}

func (s *adminUserService) GetAllUsers(page, limit int) ([]*domain.User, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.userRepo.GetAll(page, limit)
}

func (s *adminUserService) GetUserByID(id int64) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("User not found")
	}

	return user, nil
}

func (s *adminUserService) UpdateUser(user *domain.User) error {
	existing, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return errors.New("User not found")
	}

	if user.Name != "" {
		existing.Name = user.Name
	}

	if user.Email != "" {
		existing.Email = user.Email
	}

	if user.Phone != "" {
		existing.Phone = user.Phone
	}

	if user.Profile != "" {
		existing.Profile = user.Profile
	}

	if user.Role != "" {
		existing.Role = user.Role
	}

	return s.userRepo.Update(existing)
}

func (s *adminUserService) UpdateRole(id int64, role string) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	validRoles := map[string]bool{
		"admin":    true,
		"user":     true,
		"promotor": true,
	}
	if !validRoles[role] {
		return errors.New("invalid role")
	}

	return s.userRepo.UpdateRole(id, role)
}

func (s *adminUserService) DeleteUser(id int64) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}
