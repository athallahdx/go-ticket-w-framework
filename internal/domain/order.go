package domain

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID               int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           int64          `json:"user_id" gorm:"not null;index"`
	User             User           `json:"user" gorm:"foreignKey:UserID"`
	OrderItems       []OrderItem    `json:"order_items" gorm:"foreignKey:OrderID"`
	TotalAmount      float64        `json:"total_amount" gorm:"not null"`
	Status           string         `json:"status" gorm:"not null;default:'pending'"`
	PaymentMethod    string         `json:"payment_method" gorm:"not null"`
	PaymentReference string         `json:"payment_reference" gorm:"not null;uniqueIndex"`
	ExpiredAt        time.Time      `json:"expired_at" gorm:"not null"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id int64) (*Order, error)
	GetByUserID(userID int64) ([]*Order, error)
	Update(order *Order) error
	Delete(id int64) error
}

type OrderService interface {
	CreateOrder(order *Order) error
	GetOrderByID(id int64) (*Order, error)
	GetOrdersByUserID(userID int64) ([]*Order, error)
	UpdateOrder(order *Order) error
	DeleteOrder(id int64) error
}
