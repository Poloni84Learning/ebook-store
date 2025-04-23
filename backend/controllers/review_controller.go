package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewController struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewReviewController(db *gorm.DB, cfg *config.Config) *ReviewController {
	return &ReviewController{DB: db, Config: cfg}
}

// CreateReviewInput định nghĩa input cho tạo review
type CreateReviewInput struct {
	BookID  uint   `json:"book_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"max=500"`
}

// UpdateReviewInput định nghĩa input cho cập nhật review
type UpdateReviewInput struct {
	Rating  int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty,max=500"`
}
type TopRatedBook struct {
	ID            uint    `json:"id"`
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	CoverImage    string  `json:"cover_image"`
	AverageRating float64 `json:"average_rating"` // Viết hoa
}

// CreateReview tạo mới review
func (rc *ReviewController) CreateReview(c *gin.Context) {
	var input CreateReviewInput
	userID := c.GetUint("userID")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra sách tồn tại
	var book models.Book
	if err := rc.DB.First(&book, input.BookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Kiểm tra user đã review sách này chưa
	var existingReview models.Review
	if err := rc.DB.Where("user_id = ? AND book_id = ?", userID, input.BookID).First(&existingReview).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already reviewed this book"})
		return
	}

	review := models.Review{
		UserID:  userID,
		BookID:  input.BookID,
		Rating:  input.Rating,
		Comment: input.Comment,
	}

	if err := rc.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	// Cập nhật rating trung bình cho sách
	rc.updateBookAverageRating(input.BookID)
	// Tăng view count
	rc.DB.Model(&review).Update("view_count", review.ViewCount+1)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

// UpdateReview cập nhật review
func (rc *ReviewController) UpdateReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := c.GetUint("userID")
	var input UpdateReviewInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var review models.Review
	if err := rc.DB.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// Kiểm tra quyền sở hữu
	if review.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own reviews"})
		return
	}

	updates := make(map[string]interface{})
	if input.Rating != 0 {
		updates["rating"] = input.Rating
	}
	if input.Comment != "" {
		updates["comment"] = input.Comment
	}

	if err := rc.DB.Model(&review).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	// Cập nhật rating trung bình nếu thay đổi rating
	if input.Rating != 0 {
		rc.updateBookAverageRating(review.BookID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review updated successfully",
		"review":  review,
	})
}

// DeleteReview xóa review
func (rc *ReviewController) DeleteReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := c.GetUint("userID")
	var review models.Review

	if err := rc.DB.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// Kiểm tra quyền sở hữu hoặc admin
	var user models.User
	rc.DB.First(&user, userID)

	if review.UserID != userID && user.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this review"})
		return
	}

	bookID := review.BookID
	if err := rc.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	// Cập nhật rating trung bình sau khi xóa
	rc.updateBookAverageRating(bookID)

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

// GetBookReviews lấy tất cả reviews của một sách
func (rc *ReviewController) GetBookReviews(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var reviews []models.Review
	if err := rc.DB.Preload("User").
		Where("book_id = ?", bookID).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (rc *ReviewController) GetMostReviewedBooks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	period := c.DefaultQuery("period", "") // "week", "month", "year", ""

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	var timeCondition string
	switch period {
	case "week":
		timeCondition = "reviews.created_at >= DATE_TRUNC('week', NOW())"
	case "month":
		timeCondition = "reviews.created_at >= DATE_TRUNC('month', NOW())"
	case "year":
		timeCondition = "reviews.created_at >= DATE_TRUNC('year', NOW())"
	default:
		timeCondition = "1=1"
	}

	var results []struct {
		ID         uint   `json:"id"`
		Title      string `json:"title"`
		Author     string `json:"author"`
		CoverImage string `json:"cover_image"`
		ViewCount  int64  `json:"view_count"`
	}

	err = rc.DB.
		Table("reviews").
		Select("books.id, books.title, books.author, books.cover_image, COUNT(reviews.id) AS view_count").
		Joins("JOIN books ON reviews.book_id = books.id").
		Where(timeCondition).
		Group("books.id").
		Order("view_count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		log.Printf("[ERROR] Lỗi lấy most-reviewed books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy sách nhiều review nhất"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":  limit,
		"period": period,
		"data":   results,
	})
}

// GetTopRatedBooks lấy sách có điểm đánh giá trung bình cao nhất
func (rc *ReviewController) GetTopRatedBooks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	period := c.DefaultQuery("period", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	var timeCondition string
	switch period {
	case "week":
		timeCondition = "reviews.created_at >= DATE_TRUNC('week', NOW())"
	case "month":
		timeCondition = "reviews.created_at >= DATE_TRUNC('month', NOW())"
	case "year":
		timeCondition = "reviews.created_at >= DATE_TRUNC('year', NOW())"
	case "":
		timeCondition = "1=1"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Giá trị period không hợp lệ. Chọn week, month hoặc year"})
		return
	}

	var results []TopRatedBook

	err = rc.DB.
		Table("reviews").
		Select("books.id, books.title, books.author, books.cover_image, AVG(reviews.rating) as average_rating").
		Joins("JOIN books ON reviews.book_id = books.id").
		Where(timeCondition).
		Group("books.id").
		Having("COUNT(reviews.id) > 0").
		Order("average_rating DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		log.Printf("[DEBUG] Lỗi khi lấy sách rating cao nhất: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách sách top rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":  limit,
		"period": period,
		"data":   results,
	})
}

// updateBookAverageRating cập nhật rating trung bình cho sách
func (rc *ReviewController) updateBookAverageRating(bookID uint) {
	var avgRating float64
	rc.DB.Model(&models.Review{}).
		Where("book_id = ?", bookID).
		Select("COALESCE(AVG(rating), 0)").
		Row().
		Scan(&avgRating)

	rc.DB.Model(&models.Book{}).
		Where("id = ?", bookID).
		Update("average_rating", avgRating)
}
