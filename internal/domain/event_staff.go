package domain

import (
	"time"

	"gorm.io/gorm"
)

type EventStaff struct {
	ID        int64          `json:"id"         gorm:"primaryKey;autoIncrement"`
	EventID   int64          `json:"event_id"   gorm:"not null;index"`
	Event     Event          `json:"event"      gorm:"foreignKey:EventID"`
	UserID    int64          `json:"user_id"    gorm:"not null;index"`
	User      User           `json:"user"       gorm:"foreignKey:UserID"`
	Role      string         `json:"role"       gorm:"not null;default:'scanner'"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type EventStaffRepository interface {
	Create(eventStaff *EventStaff) error
	GetByID(id int64) (*EventStaff, error)
	GetByEventID(eventID int64) ([]*EventStaff, error)
	GetByUserID(userID int64) ([]*EventStaff, error)
	Update(eventStaff *EventStaff) error
	Delete(id int64) error
}

type EventStaffService interface {
	AssignStaff(eventID, userID int64, role string) error
	RevokeStaff(eventID, userID int64) error
	UpdateStaffRole(eventID, userID int64, role string) error
	GetStaffByEventID(eventID int64) ([]*EventStaff, error)
	GetStaffByUserID(userID int64) ([]*EventStaff, error)
	GetStaffRole(eventID, userID int64) (string, bool, error)
}
