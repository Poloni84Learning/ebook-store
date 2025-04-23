package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	UserID    uint   `gorm:"not null" json:"user_id"`
	BookID    uint   `gorm:"not null" json:"book_id"`
	Rating    int    `gorm:"not null" json:"rating"`
	Comment   string `gorm:"type:text" json:"comment,omitempty"`
	ViewCount int    `gorm:"default:0" json:"view_count"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user"`
	Book Book `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"book"`
}
