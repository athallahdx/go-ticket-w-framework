package domain

import (
	"time"
)

type CheckIn struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TicketID    int64     `json:"ticket_id" gorm:"not null;uniqueIndex"`
	CheckedInBy int64     `json:"checked_in_by" gorm:"not null"`
	CheckedInAt time.Time `json:"checked_in_at" gorm:"autoCreateTime"`
}

type CheckInStats struct {
	TotalTickets   int
	TotalCheckedIn int
	TotalRemaining int
	ByTicketType   map[string]int
}

type CheckInRepository interface {
	Create(checkin *CheckIn) error
	GetByID(id int64) (*CheckIn, error)
	GetByTicketID(ticketID int64) (*CheckIn, error)
	GetByEventID(eventID int64, page, limit int) ([]*CheckIn, int, error)
	GetByStaffID(staffID int64) ([]*CheckIn, error)
}

type CheckInService interface {
	CheckinTicket(ticketCode string, staffID int64) (*CheckIn, error)
	ValidateTicket(ticketCode string) (*Ticket, error)
	GetCheckinHistory(eventID int64, page, limit int) ([]*CheckIn, int, error)
	GetCheckinStats(eventID int64) (*CheckInStats, error)
}
