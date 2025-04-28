package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookController struct {
	DB     *gorm.DB
	Config *config.Config
}

type BookWithOrderCount struct {
	ID                   uint    `json:"id"`
	Title                string  `json:"title"`
	Author               string  `json:"author"`
	CoverImage           string  `json:"cover_image"`
	Price                float64 `json:"price"`
	CompletedOrdersCount int64   `json:"completed_orders_count"`
}

type BookWithoutOrderCount struct {
	ID                   uint    `json:"id"`
	Title                string  `json:"title"`
	Author               string  `json:"author"`
	CoverImage           string  `json:"cover_image"`
	Price                float64 `json:"price"`
	CompletedOrdersCount int64   `json:"-"`
}

func NewBookController(db *gorm.DB, cfg *config.Config) *BookController {
	return &BookController{DB: db, Config: cfg}
}

func (bc *BookController) CreateBook(c *gin.Context) {
	// Kiểm tra role từ JWT middleware
	userRole := c.GetString("role")
	if userRole == "" {
		log.Println("[DEBUG] Role trong context rỗng hoặc chưa set")
	} else {
		log.Printf("[DEBUG] Role trong context: %s", userRole)
	}

	if userRole != string(models.RoleAdmin) && userRole != string(models.RoleStaff) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Không có quyền thực hiện"})
		return
	}

	var input models.BookInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	book := models.Book{
		Title:       input.Title,
		Author:      input.Author,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		Category:    input.Category,
		Publisher:   input.PublishedAt,
		ISBN:        input.ISBN,
		Pages:       input.Pages,
		Language:    input.Language,
	}

	if err := bc.DB.Create(&book).Error; err != nil {
		log.Printf("[DEBUG] Không thể tạo sách: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể tạo sách"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": book})

}

func (bc *BookController) GetBooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var books []models.Book
	var total int64

	bc.DB.Model(&models.Book{}).Count(&total)
	bc.DB.Offset(offset).Limit(limit).Find(&books)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    books,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (bc *BookController) GetBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID không hợp lệ"})
		return
	}

	var book models.Book
	if err := bc.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Không tìm thấy sách"})
		return
	}
	var totalReviews int64
	err = bc.DB.Table("reviews").Where("book_id = ?", id).Count(&totalReviews).Error
	if err != nil {
		log.Printf("[DEBUG] Lỗi khi lấy tổng review: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể lấy tổng số review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": book, "reviews_count": totalReviews})
}

func (bc *BookController) UpdateBook(c *gin.Context) {
	// Kiểm tra role
	userRole := c.GetString("role")
	log.Printf("[DEBUG] Role trong context (UpdateBook): %s", userRole)
	if userRole != string(models.RoleAdmin) && userRole != string(models.RoleStaff) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Không có quyền thực hiện"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID không hợp lệ"})
		return
	}

	var input models.BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	var book models.Book
	if err := bc.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Không tìm thấy sách"})
		return
	}

	// Cập nhật thông tin
	book.Title = input.Title
	book.Author = input.Author
	book.Description = input.Description
	book.Price = input.Price
	book.Stock = input.Stock
	book.Category = input.Category
	book.Publisher = input.PublishedAt
	book.ISBN = input.ISBN
	book.Pages = input.Pages
	book.Language = input.Language

	if err := bc.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể cập nhật sách"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": book})
}

func (bc *BookController) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := bc.DB.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Không tìm thấy sách"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Lỗi khi tìm sách"})
		return
	}

	// Gọi hàm kiểm tra trước khi xoá
	if err := models.SafeDeleteBook(bc.DB, book.ID); err != nil {
		if errors.Is(err, models.ErrBookInActiveOrders) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Không thể xoá sách này vì sách đang nằm trong đơn hàng chưa xử lý!"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Xảy ra lỗi khi xoá sách: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Xoá sách thành công"})
}

func (bc *BookController) GetBookCombos(c *gin.Context) {
	bookID := c.Param("id")

	var combos []struct {
		ComboID     uint   `json:"combo_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatedBy   string `json:"created_by"`
	}

	err := bc.DB.
		Table("book_combos").
		Select("DISTINCT book_combos.id as combo_id, book_combos.title, book_combos.description, users.username as created_by").
		Joins("JOIN combo_items ON book_combos.id = combo_items.combo_id").
		Joins("JOIN users ON book_combos.created_by = users.id").
		Where("combo_items.book_id = ? AND combo_items.is_hidden = false AND book_combos.deleted_at IS NULL", bookID).
		Scan(&combos).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch combos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"book_id": bookID,
		"combos":  combos,
	})
}

func (bc *BookController) GetTopBooksByCompletedOrders(c *gin.Context) {
	period := c.DefaultQuery("period", "week") // mặc định theo tháng
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	var timeCondition string
	switch period {
	case "week":
		timeCondition = "orders.created_at  >= DATE_TRUNC('week', NOW())"
	case "month":
		timeCondition = "orders.created_at  >= DATE_TRUNC('month', NOW())"
	case "year":
		timeCondition = "orders.created_at  >= DATE_TRUNC('year', NOW())"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Giá trị period không hợp lệ. Chọn week, month hoặc year"})
		return
	}

	var results []BookWithoutOrderCount

	err = bc.DB.
		Table("order_items").
		Select("books.id, books.title, books.author, books.cover_image, books.price, SUM(order_items.quantity) as completed_orders_count").
		Joins("JOIN books ON order_items.book_id = books.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.status = ? AND "+timeCondition, "completed").
		Group("books.id").
		Order("completed_orders_count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		log.Printf("[DEBUG] Lỗi khi lấy top sách: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể lấy dữ liệu sách top order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true,
		"period": period,
		"limit":  limit,
		"data":   results,
	})
}

func (bc *BookController) GetBooksByTitle(c *gin.Context) {
	title := c.DefaultQuery("title", "")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Title is required"})
		return
	}

	var books []models.Book

	err := bc.DB.
		Where("title ILIKE ?", "%"+title+"%").
		Find(&books).Error

	if err != nil {
		log.Printf("[DEBUG] Lỗi khi lấy sách theo title: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể lấy sách theo title"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true,
		"data": books,
	})
}

func (bc *BookController) GetBooksByAuthor(c *gin.Context) {
	author := c.DefaultQuery("author", "")

	if author == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Author is required"})
		return
	}

	var books []models.Book

	err := bc.DB.
		Where("author ILIKE ?", "%"+author+"%").
		Find(&books).Error

	if err != nil {
		log.Printf("[DEBUG] Lỗi khi lấy sách theo author: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể lấy sách theo author"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true,
		"data": books,
	})
}

func (bc *BookController) GetBooksByCategory(c *gin.Context) {
	categoryParam := c.Query("category")
	if categoryParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Category is required"})
		return
	}

	category := models.BookCategory(categoryParam)
	if !category.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid category"})
		return
	}

	var books []models.Book
	if err := bc.DB.Where("category = ?", category).Find(&books).Error; err != nil {
		log.Printf("[DEBUG] Lỗi khi lấy sách theo category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể lấy sách theo category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": books})
}

func (bc *BookController) SearchBooks(c *gin.Context) {
	author := c.Query("author")
	title := c.Query("title")
	category := c.Query("category")
	description := c.Query("description")

	var books []models.Book
	query := bc.DB

	if author != "" {
		query = query.Or("author ILIKE ?", "%"+author+"%")
	}
	if title != "" {
		query = query.Or("title ILIKE ?", "%"+title+"%")
	}
	if category != "" {
		query = query.Or("category ILIKE ?", "%"+category+"%")
	}
	if description != "" {
		query = query.Or("description ILIKE ?", "%"+description+"%")
	}

	if err := query.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Lỗi khi tìm kiếm sách"})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (bc *BookController) GetTopBooks(c *gin.Context) {
	period := c.DefaultQuery("period", "yesterday") // mặc định là ngày hôm qua
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	var timeCondition string
	switch period {
	case "yesterday":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '1 day' AND orders.created_at < NOW()"
	case "day-before":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '2 day' AND orders.created_at < NOW() - INTERVAL '1 day'"
	case "last-week":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '1 week' AND orders.created_at < NOW()"
	case "last-two-weeks":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '2 weeks' AND orders.created_at < NOW() - INTERVAL '1 week'"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Giá trị period không hợp lệ"})
		return
	}

	var results []BookWithOrderCount
	err = bc.DB.
		Table("order_items").
		Select("books.id, books.title, books.author, books.cover_image, books.price, SUM(order_items.quantity) as completed_orders_count").
		Joins("JOIN books ON order_items.book_id = books.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.status = ? AND "+timeCondition, "completed").
		Group("books.id").
		Order("completed_orders_count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể lấy dữ liệu sách bán chạy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"period":  period,
		"limit":   limit,
		"data":    results,
	})
}

func (bc *BookController) GetTopCategories(c *gin.Context) {
	period := c.DefaultQuery("period", "yesterday")
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	var timeCondition string
	switch period {
	case "yesterday":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '1 day' AND orders.created_at < NOW()"
	case "day-before":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '2 day' AND orders.created_at < NOW() - INTERVAL '1 day'"
	case "last-week":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '1 week' AND orders.created_at < NOW()"
	case "last-two-weeks":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '2 weeks' AND orders.created_at < NOW() - INTERVAL '1 week'"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Giá trị period không hợp lệ"})
		return
	}

	var results []struct {
		Category             string `json:"category"`
		CompletedOrdersCount int64  `json:"completed_orders_count"`
	}

	err = bc.DB.
		Table("order_items").
		Select("books.category, SUM(order_items.quantity) as completed_orders_count").
		Joins("JOIN books ON order_items.book_id = books.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.status = ? AND "+timeCondition, "completed").
		Group("books.category").
		Order("completed_orders_count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể lấy dữ liệu category bán chạy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"period":  period,
		"limit":   limit,
		"data":    results,
	})
}

func (bc *BookController) GetTopAuthors(c *gin.Context) {
	period := c.DefaultQuery("period", "yesterday")
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	var timeCondition string
	switch period {
	case "yesterday":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '1 day' AND orders.created_at < NOW()"
	case "day-before":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '2 day' AND orders.created_at < NOW() - INTERVAL '1 day'"
	case "last-week":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '1 week' AND orders.created_at < NOW()"
	case "last-two-weeks":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '2 weeks' AND orders.created_at < NOW() - INTERVAL '1 week'"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Giá trị period không hợp lệ"})
		return
	}

	var results []struct {
		Author               string `json:"author"`
		CompletedOrdersCount int64  `json:"completed_orders_count"`
	}

	err = bc.DB.
		Table("order_items").
		Select("books.author, SUM(order_items.quantity) as completed_orders_count").
		Joins("JOIN books ON order_items.book_id = books.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.status = ? AND "+timeCondition, "completed").
		Group("books.author").
		Order("completed_orders_count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Không thể lấy dữ liệu author bán chạy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"period":  period,
		"limit":   limit,
		"data":    results,
	})
}

func (bc *BookController) GetTotalOrders(c *gin.Context) {
	period := c.DefaultQuery("period", "yesterday")

	var timeCondition string
	switch period {
	case "yesterday":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '1 day' AND orders.created_at < NOW()"
	case "day-before":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '2 day' AND orders.created_at < NOW() - INTERVAL '1 day'"
	case "last-week":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '1 week' AND orders.created_at < NOW()"
	case "last-two-weeks":
		timeCondition = "orders.created_at >= NOW() - INTERVAL '2 weeks' AND orders.created_at < NOW() - INTERVAL '1 week'"
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Giá trị period không hợp lệ"})
		return
	}

	var totalOrders int64
	err := bc.DB.Model(&models.Order{}).Where(timeCondition).Count(&totalOrders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Không thể lấy tổng số đơn hàng"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"period":       period,
		"total_orders": totalOrders,
	})
}
func (bc *BookController) GetAllCategories(c *gin.Context) {
	categories := models.GetAllBookCategories()

	var categoryStrings []string
	for _, cat := range categories {
		categoryStrings = append(categoryStrings, string(cat))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categoryStrings,
	})
}
