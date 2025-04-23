package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID        uint        `gorm:"not null" json:"user_id"`
	User          User        `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"user"`
	TotalAmount   float64     `gorm:"type:decimal(10,2);not null;check:total_amount >= 0" json:"total_amount"`
	Status        string      `gorm:"type:varchar(20);default:'pending'" json:"status"`
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"` // Liên kết với OrderItem
	PaymentMethod string      `gorm:"type:varchar(20);default:'Card'" json:"payment_method"`
}
