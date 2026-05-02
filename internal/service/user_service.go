package service

import (
	"go-ticket/internal/domain"
	"go-ticket/pkg/file"
	"mime/multipart"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetProfileByID(id int64) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateProfile(id int64, input domain.UpdateProfileInput, fileHeader *multipart.FileHeader) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	finalName := user.Name
	if input.Name != nil && *input.Name != "" {
		finalName = *input.Name
	}

	if fileHeader != nil {
		newPath, err := file.ReplaceFile(
			user.Profile,
			fileHeader,
			file.EntityUser,
			user.ID,
			finalName,
		)
		if err != nil {
			log.Error().Err(err).Int64("user_id", id).Msg("Failed to replace profile image file")
			return nil, err
		}

		user.Profile = filepath.ToSlash(newPath)
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Phone != nil {
		user.Phone = *input.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		log.Error().Err(err).Int64("user_id", id).Msg("Failed to update user profile in database")
		return nil, err
	}

	return user, nil
}
