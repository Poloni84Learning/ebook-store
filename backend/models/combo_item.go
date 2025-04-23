package models

import "gorm.io/gorm"

type ComboItem struct {
	gorm.Model
	ComboID uint `gorm:"not null;index:idx_combo_book,unique,where:deleted_at is null"`
	BookID  uint `gorm:"not null;index:idx_combo_book,unique,where:deleted_at is null"`

	IsHidden bool `gorm:"default:false"` // Thêm trường mới
	Book     Book `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
