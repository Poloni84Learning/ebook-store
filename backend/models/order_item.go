package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	OrderID  uint    `gorm:"not null;index"`                                                                // Liên kết với Order
	BookID   uint    `gorm:"not null;index"`                                                                // Liên kết với Book
	Quantity int     `gorm:"not null;check:quantity > 0"`                                                   // Số lượng sách trong đơn hàng
	Price    float64 `gorm:"type:decimal(10,2);not null"`                                                   // Giá của sách
	Discount float64 `gorm:"type:decimal(10,2);default:0.00"`                                               // Giảm giá cho sản phẩm, nếu có
	Book     Book    `gorm:"foreignKey:BookID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // Liên kết với Book
}

type OrderItemResponse struct {
	ID       uint    `json:"id"`
	BookID   uint    `json:"book_id"`
	Title    string  `json:"title,omitempty"`
	Author   string  `json:"author,omitempty"`    // Thêm tác giả
	ImageURL string  `json:"image_url,omitempty"` // Thêm ảnh cover
	Quantity int     `json:"quantity"`
	Discount float64 `json:"discount"`
	Price    float64 `json:"price"` // Giá
}

func (oi *OrderItem) ToResponse() OrderItemResponse {
	// Khởi tạo giá trị mặc định
	resp := OrderItemResponse{
		ID:       oi.ID,
		BookID:   oi.BookID,
		Quantity: oi.Quantity,
		Discount: oi.Discount,
		Price:    oi.Price,
	}

	// Kiểm tra quan hệ Book đã được preload chưa thông qua BookID
	if oi.Book.ID != 0 { // Nếu Book đã được preload (ID khác 0)
		resp.Title = oi.Book.Title
		resp.Author = oi.Book.Author
		resp.ImageURL = oi.Book.CoverImage
	}

	return resp
}
