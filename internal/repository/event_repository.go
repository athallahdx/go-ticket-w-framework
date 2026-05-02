package repository

import (
	"errors"
	"go-ticket/internal/domain"

	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) domain.EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) Create(event *domain.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) GetByID(id int64) (*domain.Event, error) {
	var event domain.Event
	err := r.db.Preload("Organizer").Preload("Images").Preload("TicketTypes").First(&event, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetAll(page, limit int) ([]*domain.Event, int, error) {
	var events []*domain.Event
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&domain.Event{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Organizer").Offset(offset).Limit(limit).Order("created_at DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}
	return events, int(total), nil
}

func (r *eventRepository) Update(event *domain.Event) error {
	result := r.db.Model(event).Updates(event)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (r *eventRepository) Delete(id int64) error {
	result := r.db.Delete(&domain.Event{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("event not found or already deleted")
	}
	return nil
}
