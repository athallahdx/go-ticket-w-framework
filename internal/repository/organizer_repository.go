package repository

import (
	"errors"
	"go-ticket/internal/domain"

	"gorm.io/gorm"
)

type organizerRepository struct {
	db *gorm.DB
}

func NewOrganizerRepository(db *gorm.DB) domain.OrganizerRepository {
	return &organizerRepository{
		db: db,
	}
}

func (r *organizerRepository) Create(organizer *domain.Organizer) error {
	return r.db.Create(organizer).Error
}

func (r *organizerRepository) GetByID(id int64) (*domain.Organizer, error) {
	var organizer domain.Organizer
	err := r.db.Preload("User").First(&organizer, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organizer not found")
		}
		return nil, err
	}
	return &organizer, nil
}

func (r *organizerRepository) GetByUserID(userID int64) (*domain.Organizer, error) {
	var organizer domain.Organizer
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&organizer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organizer not found")
		}
		return nil, err
	}
	return &organizer, nil
}

func (r *organizerRepository) GetAll(page, limit int) ([]*domain.Organizer, int, error) {
	var organizers []*domain.Organizer
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&domain.Organizer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("User").Offset(offset).Limit(limit).Order("created_at DESC").Find(&organizers).Error
	if err != nil {
		return nil, 0, err
	}
	return organizers, int(total), nil
}

func (r *organizerRepository) Update(organizer *domain.Organizer) error {
	result := r.db.Model(organizer).Updates(organizer)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("organizer not found")
	}
	return nil
}

func (r *organizerRepository) Delete(id int64) error {
	result := r.db.Delete(&domain.Organizer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("organizer not found or already deleted")
	}
	return nil
}
