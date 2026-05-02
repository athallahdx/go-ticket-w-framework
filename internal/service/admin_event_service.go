package service

import (
	"errors"
	"go-ticket/internal/domain"
)

type adminEventService struct {
	eventRepo domain.EventRepository
}

func NewAdminEventService(eventRepo domain.EventRepository) domain.AdminEventService {
	return &adminEventService{
		eventRepo: eventRepo,
	}
}

func (s *adminEventService) CreateEvent(event *domain.Event) error {
	return s.eventRepo.Create(event)
}

func (s *adminEventService) GetAllEvents(page, limit int) ([]*domain.Event, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.eventRepo.GetAll(page, limit)
}

func (s *adminEventService) GetEventByID(id int64) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("Event not found")
	}
	return event, nil
}

func (s *adminEventService) UpdateEvent(event *domain.Event) error {
	existing, err := s.eventRepo.GetByID(event.ID)
	if err != nil {
		return errors.New("Event not found")
	}

	if event.Name != "" {
		existing.Name = event.Name
	}
	if event.Description != "" {
		existing.Description = event.Description
	}
	if event.Location != "" {
		existing.Location = event.Location
	}
	if event.Thumbnail != "" {
		existing.Thumbnail = event.Thumbnail
	}
	if !event.Date.IsZero() {
		existing.Date = event.Date
	}

	return s.eventRepo.Update(existing)
}

func (s *adminEventService) DeleteEvent(id int64) error {
	_, err := s.eventRepo.GetByID(id)
	if err != nil {
		return errors.New("event not found")
	}
	return s.eventRepo.Delete(id)
}
