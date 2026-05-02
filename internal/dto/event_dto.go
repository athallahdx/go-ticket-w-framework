package dto

import "time"

type CreateEventRequest struct {
	OrganizerID int64     `json:"organizer_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Thumbnail   string    `json:"thumbnail"`
	Date        time.Time `json:"date" binding:"required"`
}

type UpdateEventRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Thumbnail   string    `json:"thumbnail"`
	Date        time.Time `json:"date"`
}

type EventResponse struct {
	ID          int64     `json:"id"`
	OrganizerID int64     `json:"organizer_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Thumbnail   string    `json:"thumbnail"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventListResponse struct {
	Events []EventResponse `json:"events"`
	Total  int64           `json:"total"`
}
