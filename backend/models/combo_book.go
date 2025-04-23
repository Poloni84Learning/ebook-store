package models

import (
	"gorm.io/gorm"
)

type BookCombo struct {
	gorm.Model
	Title       string      `gorm:"size:100;not null;index;check:title <> ''"`
	Description string      `gorm:"type:text;not null"`
	CreatedBy   uint        `gorm:"not null"`
	User        User        `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComboItems  []ComboItem `gorm:"foreignKey:ComboID"`
}
