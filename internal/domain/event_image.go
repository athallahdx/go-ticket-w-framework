package domain

import "time"

type EventImage struct {
	ID        int64     `json:"id"             gorm:"primaryKey;autoIncrement"`
	EventID   int64     `json:"event_id"       gorm:"not null;index"`
	ImageURL  string    `json:"image_url"      gorm:"not null"`
	CreatedAt time.Time `json:"created_at"     gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at"     gorm:"autoUpdateTime"`
}
