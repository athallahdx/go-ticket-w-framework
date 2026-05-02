package domain

import (
	"time"
)

type OrderItem struct {
	ID           int64      `json:"id"             gorm:"primaryKey;autoIncrement"`
	OrderID      int64      `json:"order_id"       gorm:"not null;index"`
	TicketTypeID int64      `json:"ticket_type_id" gorm:"not null;index"`
	TicketType   TicketType `json:"ticket_type"    gorm:"foreignKey:TicketTypeID"`
	Quantity     int        `json:"quantity"       gorm:"not null"`
	Price        float64    `json:"price"          gorm:"not null"`
	CreatedAt    time.Time  `json:"created_at"     gorm:"autoCreateTime"`
}

type OrderItemRepository interface {
	Create(orderItem *OrderItem) error
	GetByID(id int64) (*OrderItem, error)
	GetByOrderID(orderID int64) ([]*OrderItem, error)
	GetByTicketTypeID(ticketTypeID int64) ([]*OrderItem, error)
	Update(orderItem *OrderItem) error
	Delete(id int64) error
}

type OrderItemService interface {
	CreateOrderItem(orderItem *OrderItem) error
	GetOrderItemByID(id int64) (*OrderItem, error)
	GetOrderItemsByOrderID(orderID int64) ([]*OrderItem, error)
	GetOrderItemsByTicketTypeID(ticketTypeID int64) ([]*OrderItem, error)
	UpdateOrderItem(orderItem *OrderItem) error
	DeleteOrderItem(id int64) error
}
