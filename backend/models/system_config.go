package models

import "gorm.io/gorm"

type SystemConfig struct {
	gorm.Model
	ShippingFee   float64 `gorm:"type:decimal(10,2);default:0.0;not null"`
	Promotion     float64 `gorm:"type:decimal(5,2);default:0;not null"`
	PromotionInfo string  `gorm:"type:text"`
}
