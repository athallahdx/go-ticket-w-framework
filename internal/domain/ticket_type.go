package domain

import (
	"time"

	"gorm.io/gorm"
)

type TicketType struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID     int64          `json:"event_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`
	Price       float64        `json:"price" gorm:"not null"`
	Quota       int            `json:"quota" gorm:"not null"`
	Sold        int            `json:"sold" gorm:"not null;default:0"`
	Description string         `json:"description" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type TicketTypeRepository interface {
	Create(ticketType *TicketType) error
	GetByID(id int64) (*TicketType, error)
	GetByEventID(eventID int64) ([]*TicketType, error)
	Update(ticketType *TicketType) error
	Delete(id int64) error
}

type TicketTypeService interface {
	CreateTicketType(ticketType *TicketType) error
	UpdateTicketType(ticketType *TicketType) error
	DeleteTicketType(id int64) error
	GetTicketTypeByID(id int64) (*TicketType, error)
	GetTicketTypesByEventID(eventID int64) ([]*TicketType, error)
	GetAvailableTicketTypes(eventID int64) ([]*TicketType, error)
}
