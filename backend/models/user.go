package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleStaff    Role = "staff"
	RoleCustomer Role = "customer"
)

type User struct {
	gorm.Model
	Username     string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string     `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	FirstName    string     `gorm:"size:50" json:"first_name,omitempty"`
	LastName     string     `gorm:"size:50" json:"last_name,omitempty"`
	Phone        string     `gorm:"size:20;index" json:"phone,omitempty"` // Thêm index cho tìm kiếm
	Address      string     `gorm:"size:255" json:"address,omitempty"`
	Role         Role       `gorm:"type:varchar(20);default:'customer';index" json:"role"` // Thêm index
	LastLogin    *time.Time `json:"last_login,omitempty"`
	IsActive     bool       `gorm:"default:true;index" json:"is_active"` // Thêm index
	AvatarURL    string     `gorm:"size:255" json:"avatar_url,omitempty"`

	// Quan hệ
	Orders  []Order  `gorm:"foreignKey:UserID" json:"orders,omitempty"` // Nên bật omitempty
	Reviews []Review `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
}

// BeforeCreate Hook: đảm bảo user có Role
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = RoleCustomer
	}
	return nil
}

// ToResponse: trả về định dạng public cho API (không có mật khẩu)
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: &u.FirstName,
		LastName:  &u.LastName,
		Role:      u.Role,
		AvatarURL: &u.AvatarURL,
		Address:   &u.Address,
		Phone:     &u.Phone,
	}
}

// UserResponse: cấu trúc dữ liệu trả ra cho client
type UserResponse struct {
	ID        uint    `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Role      Role    `json:"role"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Address   *string `json:"address,omitempty"`
	Phone     *string `json:"phone,omitempty"`
}

// IsAdmin: check quyền admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsStaff: check quyền staff hoặc admin
func (u *User) IsStaff() bool {
	return u.Role == RoleStaff || u.Role == RoleAdmin
}
