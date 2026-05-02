package dto

import "time"

type CreateOrganizerRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	CompanyName string `json:"company_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	City        string `json:"city" binding:"required"`
	Province    string `json:"province" binding:"required"`
}

type UpdateOrganizerRequest struct {
	CompanyName string `json:"company_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	City        string `json:"city"`
	Province    string `json:"province"`
}

type OrganizerResponse struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	CompanyName string     `json:"company_name"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	Logo        string     `json:"logo"`
	Description string     `json:"description"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	IsVerified  bool       `json:"is_verified"`
	VerifiedAt  *time.Time `json:"verified_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type OrganizerListResponse struct {
	Organizers []OrganizerResponse `json:"organizers"`
	Total      int64               `json:"total"`
}
