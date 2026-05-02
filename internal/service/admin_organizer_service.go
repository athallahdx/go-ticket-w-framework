package service

import (
	"errors"
	"go-ticket/internal/domain"
	"go-ticket/pkg/file"
	"mime/multipart"
)

type adminOrganizerService struct {
	organizerRepo domain.OrganizerRepository
}

func NewAdminOrganizerService(organizerRepo domain.OrganizerRepository) domain.AdminOrganizerService {
	return &adminOrganizerService{
		organizerRepo: organizerRepo,
	}
}

func (s *adminOrganizerService) CreateOrganizer(organizer *domain.Organizer, fileHeader *multipart.FileHeader) error {
	// check if user already has an organizer profile
	existing, _ := s.organizerRepo.GetByUserID(organizer.UserID)
	if existing != nil {
		return errors.New("user already has an organizer profile")
	}

	if err := s.organizerRepo.Create(organizer); err != nil {
		return err
	}

	if fileHeader != nil {
		logoPath, err := file.SaveSingleFile(fileHeader, file.EntityOrganizer, organizer.ID, organizer.CompanyName)
		if err != nil {
			s.organizerRepo.Delete(organizer.ID) // cleanup
			return err
		}
		organizer.Logo = logoPath
		s.organizerRepo.Update(organizer)
	}

	return nil
}

func (s *adminOrganizerService) GetAllOrganizers(page, limit int) ([]*domain.Organizer, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.organizerRepo.GetAll(page, limit)
}

func (s *adminOrganizerService) GetOrganizerByID(id int64) (*domain.Organizer, error) {
	organizer, err := s.organizerRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("Organizer not found")
	}
	return organizer, nil
}

func (s *adminOrganizerService) UpdateOrganizer(organizer *domain.Organizer, fileHeader *multipart.FileHeader) error {
	existing, err := s.organizerRepo.GetByID(organizer.ID)
	if err != nil {
		return errors.New("Organizer not found")
	}

	if organizer.CompanyName != "" {
		existing.CompanyName = organizer.CompanyName
	}
	if organizer.Phone != "" {
		existing.Phone = organizer.Phone
	}
	if organizer.Email != "" {
		existing.Email = organizer.Email
	}
	if organizer.Description != "" {
		existing.Description = organizer.Description
	}
	if organizer.City != "" {
		existing.City = organizer.City
	}
	if organizer.Province != "" {
		existing.Province = organizer.Province
	}

	if fileHeader != nil {
		logoPath, err := file.ReplaceFile(existing.Logo, fileHeader, file.EntityOrganizer, existing.ID, existing.CompanyName)
		if err != nil {
			return err
		}
		existing.Logo = logoPath
	}

	return s.organizerRepo.Update(existing)
}

func (s *adminOrganizerService) DeleteOrganizer(id int64) error {
	_, err := s.organizerRepo.GetByID(id)
	if err != nil {
		return errors.New("organizer not found")
	}
	return s.organizerRepo.Delete(id)
}
