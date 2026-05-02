package domain

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	ID           int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Code         string         `json:"code" gorm:"not null;uniqueIndex"`
	UserID       int64          `json:"user_id" gorm:"not null;index"`
	User         User           `json:"user" gorm:"foreignKey:UserID"`
	OrderID      int64          `json:"order_id" gorm:"not null;index"`
	TicketTypeID int64          `json:"ticket_type_id" gorm:"not null;index"`
	TicketType   TicketType     `json:"ticket_type" gorm:"foreignKey:TicketTypeID"`
	QRCode       string         `json:"qr_code" gorm:"not null"`
	Status       string         `json:"status" gorm:"not null;default:'active'"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type TicketRepository interface {
	Create(ticket *Ticket) error
	GetByID(id int64) (*Ticket, error)
	GetByCode(code string) (*Ticket, error)
	GetByUserID(userID int64) ([]*Ticket, error)
	GetByTicketTypeID(ticketTypeID int64) ([]*Ticket, error)
	Update(ticket *Ticket) error
	Delete(id int64) error
}

type TicketService interface {
	CreateTicket(ticket *Ticket) error
	GetTicketByID(id int64) (*Ticket, error)
	GetTicketByCode(code string) (*Ticket, error)
	GetTicketsByUserID(userID int64) ([]*Ticket, error)
	GetTicketsByTicketTypeID(ticketTypeID int64) ([]*Ticket, error)
	UpdateTicket(ticket *Ticket) error
	DeleteTicket(id int64) error
}
