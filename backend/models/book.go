package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BookCategory string

const (
	CategoryFiction     BookCategory = "Fiction"
	CategoryNonFiction  BookCategory = "Non-fiction"
	CategoryScience     BookCategory = "Science"
	CategoryMath        BookCategory = "Math"
	CategoryTechnology  BookCategory = "Technology"
	CategoryHistory     BookCategory = "History"
	CategoryBiography   BookCategory = "Biography"
	CategoryPhilosophy  BookCategory = "Philosophy"
	CategorySelfHelp    BookCategory = "Self-help"
	CategoryChildren    BookCategory = "Children"
	CategoryEducation   BookCategory = "Education"
	CategoryComics      BookCategory = "Comics"
	CategoryFantasy     BookCategory = "Fantasy"
	CategoryMystery     BookCategory = "Mystery"
	CategoryHorror      BookCategory = "Horror"
	CategoryRomance     BookCategory = "Romance"
	CategoryBusiness    BookCategory = "Business"
	CategoryProgramming BookCategory = "Programming"
)

type Book struct {
	gorm.Model
	Title         string       `gorm:"size:100;not null;index" json:"title"`
	Author        string       `gorm:"size:50;not null;index" json:"author"`
	Description   string       `gorm:"type:text" json:"description,omitempty"`
	Price         float64      `gorm:"type:decimal(10,2);not null;check:price > 0" json:"price"`
	Stock         int          `gorm:"default:0;not null;check:stock >= 0" json:"stock"`
	CoverImage    string       `gorm:"size:255" json:"cover_image,omitempty"`
	Category      BookCategory `gorm:"size:50;index" json:"category"`
	Publisher     string       `gorm:"size:50" json:"publisher,omitempty"`
	ISBN          string       `gorm:"size:20;uniqueIndex;not null" json:"isbn"`
	Pages         int          `gorm:"check:pages >= 1" json:"pages,omitempty"`
	Language      string       `gorm:"size:20" json:"language,omitempty"`
	PublishedAt   string       `json:"published_at"`
	AverageRating float64      `gorm:"type:decimal(3,2);default:0" json:"average_rating"`

	OrderItems []OrderItem `gorm:"foreignKey:BookID" json:"-"`
	Reviews    []Review    `gorm:"foreignKey:BookID" json:"-"`
}

type BookInput struct {
	Title       string       `json:"title" binding:"required,min=3,max=100"`
	Author      string       `json:"author" binding:"required,min=3,max=50"`
	Description string       `json:"description" binding:"max=500"`
	Price       float64      `json:"price" binding:"required,gt=0"`
	Stock       int          `json:"stock" binding:"gte=0"`
	CoverImage  string       `json:"cover_image"`
	Category    BookCategory `json:"category" binding:"max=50"`
	PublishedAt string       `json:"published_at" binding:"max=50"`
	ISBN        string       `json:"isbn" binding:"required,min=10,max=20"`
	Pages       int          `json:"pages" binding:"gte=1"`
	Language    string       `json:"language" binding:"max=20"`
}

type BookResponse struct {
	ID            uint         `json:"id"`
	Title         string       `json:"title"`
	Author        string       `json:"author"`
	Description   string       `json:"description,omitempty"`
	Price         float64      `json:"price"`
	Stock         int          `json:"stock"`
	CoverImage    string       `json:"cover_image,omitempty"`
	Category      BookCategory `json:"category"`
	PublishedAt   string       `json:"published_at,omitempty"`
	ISBN          string       `json:"isbn"`
	Pages         int          `json:"pages,omitempty"`
	Language      string       `json:"language,omitempty"`
	AverageRating float64      `json:"average_rating"`
	CreatedAt     time.Time    `json:"created_at"`
}

func (b *Book) ToResponse() *BookResponse {
	return &BookResponse{
		ID:            b.ID,
		Title:         b.Title,
		Author:        b.Author,
		Description:   b.Description,
		Price:         b.Price,
		Stock:         b.Stock,
		CoverImage:    b.CoverImage,
		Category:      b.Category,
		PublishedAt:   b.PublishedAt,
		ISBN:          b.ISBN,
		Pages:         b.Pages,
		Language:      b.Language,
		AverageRating: b.AverageRating,
		CreatedAt:     b.CreatedAt,
	}
}

func (b *Book) BeforeCreate(tx *gorm.DB) error {
	if !b.Category.IsValid() {
		return fmt.Errorf("invalid category: %s", b.Category)
	}
	if b.Stock < 0 {
		b.Stock = 0
	}
	if b.Price <= 0 {
		b.Price = 1.0
	}
	if b.CoverImage == "" {
		b.CoverImage = "/uploads/default-cover.jpg"
	}
	return nil
}

var ErrBookInActiveOrders = errors.New("không thể xoá sách này vì đang có đơn hàng chưa xử lý")

// Hàm kiểm tra & xoá an toàn Book
func SafeDeleteBook(db *gorm.DB, bookID uint) error {
	var count int64

	// Kiểm tra xem sách có nằm trong Order chưa xử lý không
	err := db.Model(&OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("order_items.book_id = ? AND orders.status IN ?", bookID, []string{"pending", "confirmed"}).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return ErrBookInActiveOrders
	}
	// Thêm bước cập nhật is_hidden trước khi xóa
	if err := db.Model(&ComboItem{}).Where("book_id = ?", bookID).Update("is_hidden", true).Error; err != nil {
		return err
	}
	// Cho phép xóa mềm (soft delete)
	if err := db.Delete(&Book{}, bookID).Error; err != nil {
		return err
	}

	return nil
}

var validCategories = map[BookCategory]bool{
	CategoryFiction:     true,
	CategoryNonFiction:  true,
	CategoryScience:     true,
	CategoryMath:        true,
	CategoryTechnology:  true,
	CategoryHistory:     true,
	CategoryBiography:   true,
	CategoryPhilosophy:  true,
	CategorySelfHelp:    true,
	CategoryChildren:    true,
	CategoryEducation:   true,
	CategoryComics:      true,
	CategoryFantasy:     true,
	CategoryMystery:     true,
	CategoryHorror:      true,
	CategoryRomance:     true,
	CategoryBusiness:    true,
	CategoryProgramming: true,
}

func (c BookCategory) IsValid() bool {
	return validCategories[c]
}
