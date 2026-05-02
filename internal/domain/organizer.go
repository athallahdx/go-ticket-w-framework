package domain

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Organizer struct {
	ID          int64          `json:"id"           gorm:"primaryKey;autoIncrement"`
	UserID      int64          `json:"user_id"      gorm:"not null;index"`
	User        User           `json:"user"         gorm:"foreignKey:UserID;references:ID"`
	CompanyName string         `json:"company_name" gorm:"not null"`
	Phone       string         `json:"phone"        gorm:"not null"`
	Email       string         `json:"email"        gorm:"not null;uniqueIndex"`
	Logo        string         `json:"logo"         gorm:"not null"`
	Description string         `json:"description"  gorm:"not null"`
	City        string         `json:"city"         gorm:"not null"`
	Province    string         `json:"province"     gorm:"not null"`
	IsVerified  bool           `json:"is_verified"  gorm:"not null;default:false"`
	VerifiedAt  *time.Time     `json:"verified_at"`
	VerifiedBy  *int64         `json:"verified_by"`
	CreatedAt   time.Time      `json:"created_at"   gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at"   gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"   gorm:"index"`
}

type OrganizerRepository interface {
	Create(organizer *Organizer) error
	GetByID(id int64) (*Organizer, error)
	GetByUserID(userID int64) (*Organizer, error)
	GetAll(page, limit int) ([]*Organizer, int, error)
	Update(organizer *Organizer) error
	Delete(id int64) error
}

type OrganizerService interface {
	RegisterAsOrganizer(user *User, companyName string) (*Organizer, error)
	GetOrganizerByID(id int64) (*Organizer, error)
	GetOrganizerByUserID(userID int64) (*Organizer, error)
	GetAllOrganizers(page, limit int) ([]*Organizer, int, error)
	UpdateOrganizer(organizer *Organizer) error
	DeleteOrganizer(id int64) error

	VerifyOrganizer(organizerID int64, adminID int64) error
	RevokeVerification(organizerID int64, adminID int64) error
}

type AdminOrganizerService interface {
	CreateOrganizer(organizer *Organizer, fileHeader *multipart.FileHeader) error
	GetAllOrganizers(page, limit int) ([]*Organizer, int, error)
	GetOrganizerByID(id int64) (*Organizer, error)
	UpdateOrganizer(organizer *Organizer, fileHeader *multipart.FileHeader) error
	DeleteOrganizer(id int64) error
}
