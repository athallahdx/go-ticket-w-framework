package domain

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID          int64          `json:"id"           gorm:"primaryKey;autoIncrement"`
	OrganizerID int64          `json:"organizer_id" gorm:"not null;index"`
	Organizer   Organizer      `json:"organizer"    gorm:"foreignKey:OrganizerID"`
	Thumbnail   string         `json:"thumbnail"    gorm:"not null"`
	Images      []EventImage   `json:"images"       gorm:"foreignKey:EventID"`
	TicketTypes []TicketType   `json:"ticket_types" gorm:"foreignKey:EventID"`
	Name        string         `json:"name"         gorm:"not null"`
	Description string         `json:"description"  gorm:"not null"`
	Location    string         `json:"location"     gorm:"not null"`
	Date        time.Time      `json:"date"         gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"   gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at"   gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"   gorm:"index"`
}

type EventRepository interface {
	Create(event *Event) error
	GetByID(id int64) (*Event, error)
	GetAll(page, limit int) ([]*Event, int, error)
	Update(event *Event) error
	Delete(id int64) error
}

type EventFilter struct {
	OrganizerID int64
	City        string
	DateFrom    time.Time
	DateTo      time.Time
	Search      string
}

type EventService interface {
	CreateEvent(event *Event) error
	UpdateEvent(event *Event) error
	DeleteEvent(id int64) error
	GetEventByID(id int64) (*Event, error)
	GetAllEvents(filter EventFilter, page, limit int) ([]*Event, int, error)
	GetEventsByOrganizerID(organizerID int64, page, limit int) ([]*Event, int, error)

	AddImage(eventID int64, imageURL string) (*EventImage, error)
	RemoveImage(imageID int64) error
	GetImagesByEventID(eventID int64) ([]*EventImage, error)
}

type AdminEventService interface {
	CreateEvent(event *Event) error
	GetAllEvents(page, limit int) ([]*Event, int, error)
	GetEventByID(id int64) (*Event, error)
	UpdateEvent(event *Event) error
	DeleteEvent(id int64) error
}
