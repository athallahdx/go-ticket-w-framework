package domain

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"not null;uniqueIndex"`
	Phone     string         `json:"phone"`
	Profile   string         `json:"profile"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"not null;default:'user'"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id int64) (*User, error)
	UpdateRole(id int64, role string) error
	Delete(id int64) error
	Update(user *User) error
	GetAll(page, limit int) ([]*User, int, error)
	GetAllWithDeleted(page, limit int) ([]*User, int, error)
}

type AdminUserService interface {
	GetAllUsers(page, limit int) ([]*User, int, error)
	GetUserByID(id int64) (*User, error)
	UpdateUser(user *User) error
	UpdateRole(id int64, role string) error
	DeleteUser(id int64) error
}

type UserService interface {
	GetProfileByID(id int64) (*User, error)
	UpdateProfile(id int64, input UpdateProfileInput, fileHeader *multipart.FileHeader) (*User, error)
}
