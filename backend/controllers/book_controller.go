package controllers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type BookController struct {
	DB         *gorm.DB
	Config     *config.Config
	tempTokens map[string]string // Dùng cho demo
	sync.Mutex
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

type PDFExtractionResult struct {
	Keywords  []string `json:"keywords"`
	TocTitles []string `json:"toc_titles"`
	Status    string   `json:"status"`
	Message   string   `json:"message,omitempty"`
}

func NewBookController(db *gorm.DB, cfg *config.Config) *BookController {
	return &BookController{DB: db, Config: cfg, tempTokens: make(map[string]string)}
}

func (bc *BookController) CreateBook(c *gin.Context) {
	// 1. Kiểm tra quyền
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

	// 2. Parse form data với giới hạn kích thước file
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể parse form data", "details": err.Error()})
		return
	}

	// 3. Lấy các giá trị từ form
	formValues := map[string]string{
		"title":        c.PostForm("title"),
		"author":       c.PostForm("author"),
		"description":  c.PostForm("description"),
		"price":        c.PostForm("price"),
		"stock":        c.PostForm("stock"),
		"category":     c.PostForm("category"),
		"publisher":    c.PostForm("publisher"),
		"isbn":         c.PostForm("isbn"),
		"pages":        c.PostForm("pages"),
		"language":     c.PostForm("language"),
		"published_at": c.PostForm("published_at"),
		"toc_pages":    c.PostForm("toc_pages"), // Thay đổi từ "toc" thành "toc_pages"
	}

	// 4. Validate các trường bắt buộc
	// for field, value := range formValues {
	// 	if value == "" && field != "published_at" { // published_at có thể optional
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Thiếu trường bắt buộc: %s", field)})
	// 		return
	// 	}
	// }

	// 5. Convert kiểu dữ liệu
	price, err := strconv.ParseFloat(formValues["price"], 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Giá tiền không hợp lệ"})
		return
	}

	stock, err := strconv.Atoi(formValues["stock"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Số lượng không hợp lệ"})
		return
	}

	pages, err := strconv.Atoi(formValues["pages"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Số trang không hợp lệ"})
		return
	}

	// 6. Xử lý file PDF
	pdfFile, err := c.FormFile("pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu file PDF", "details": err.Error()})
		return
	}
	var imageUrl string
	if imageFile, err := c.FormFile("cover_image"); err == nil {
		// Lưu file ảnh
		imageUrl, err = SaveImageFile(imageFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu ảnh", "details": err.Error()})
			return
		}
	}

	// 7. Lưu file PDF tạm để xử lý
	tempPDFPath, err := saveTempFile(pdfFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu file tạm", "details": err.Error()})
		return
	}
	defer os.Remove(tempPDFPath) // Xóa file tạm sau khi xử lý xong

	// 10. Lưu file PDF vĩnh viễn
	pdfUrl, err := SavePDFFile(pdfFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu file PDF", "details": err.Error()})
		return
	}

	// --- Gọi API để lấy keywords và toc_titles ---
	keywords, tocTitles, err := callKeywordAndTOCApi(tempPDFPath, formValues["title"], formValues["author"], string(formValues["category"]), formValues["toc_pages"])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy dữ liệu phân tích PDF", "details": err.Error()})
		return
	}
	// 11. Tạo book record
	book := models.Book{
		Title:       formValues["title"],
		Author:      formValues["author"],
		Description: formValues["description"],
		Price:       price,
		Stock:       stock,
		Category:    models.BookCategory(formValues["category"]),
		Publisher:   formValues["publisher"],
		ISBN:        formValues["isbn"],
		Pages:       pages,
		Language:    formValues["language"],
		PublishedAt: formValues["published_at"],
		PDFUrl:      pdfUrl,
		CoverImage:  imageUrl,
		Keywords:    keywords,
		TOCTitles:   tocTitles,
	}

	if err := bc.DB.Create(&book).Error; err != nil {
		log.Printf("[ERROR] Failed to create book: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Không thể tạo sách",
			"details": err.Error(),
		})
		return
	}

	// 12. Trả về response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"book":       book,
			"keywords":   keywords,
			"toc_titles": tocTitles,
		},
	})
}

func callKeywordAndTOCApi(pdfPath, title, author, topic, tocPages string) (pq.StringArray, pq.StringArray, error) {

	file, err := os.Open(pdfPath)
	if err != nil {
		return nil, nil, fmt.Errorf("mở file PDF thất bại: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// file
	part, err := writer.CreateFormFile("file", filepath.Base(pdfPath))
	if err != nil {
		return nil, nil, fmt.Errorf("tạo form file thất bại: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, nil, fmt.Errorf("ghi dữ liệu file thất bại: %w", err)
	}

	// required fields
	_ = writer.WriteField("book_title", title)
	_ = writer.WriteField("authors", author)
	_ = writer.WriteField("topic", topic)

	// optional
	if tocPages != "" {
		_ = writer.WriteField("toc_pages", tocPages)
	}

	if err := writer.Close(); err != nil {
		return nil, nil, fmt.Errorf("đóng writer thất bại: %w", err)
	}

	// Gửi request
	apiURL := fmt.Sprintf("%s/extract-keywords", getAPIBaseURL())
	req, err := http.NewRequest("POST", apiURL, body) // chỉnh lại URL phù hợp
	if err != nil {
		return nil, nil, fmt.Errorf("tạo request thất bại: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 4. Cấu hình HTTP client
	client := &http.Client{
		Timeout: 0, // Vô hiệu hóa timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("gửi request thất bại: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("API trả về lỗi %d: %s", resp.StatusCode, string(responseBody))
	}

	var result struct {
		Keywords  []string `json:"keywords"`
		TOCTitles []string `json:"toc_titles"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("decode JSON thất bại: %w", err)
	}
	keywords := pq.StringArray(result.Keywords)
	tocTitles := pq.StringArray(result.TOCTitles)
	return keywords, tocTitles, nil
}

func getAPIBaseURL() string {
	// Khi chạy trong Docker (tự động phát hiện)
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "http://host.docker.internal:8001" // Cho Windows/Mac Docker
	}
	return "http://localhost:8001" // Khi chạy trực tiếp
}

// Helper function để lưu file tạm
func saveTempFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, src); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func (bc *BookController) GenerateDownloadLink(c *gin.Context) {
	// Lấy bookID từ URL params
	bookIDStr := c.Param("id")
	var bookID uint
	if _, err := fmt.Sscanf(bookIDStr, "%d", &bookID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sách không hợp lệ"})
		return
	}

	// Tạo token ngẫu nhiên
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}
	token := hex.EncodeToString(tokenBytes)

	// Lưu token vào map với mutex
	bc.Lock()
	bc.tempTokens[token] = bookIDStr // Lưu dưới dạng string để tiện xử lý
	bc.Unlock()

	// Tạo URL tải về (sử dụng config nếu có, hoặc localhost:8081)
	baseURL := "http://localhost:8081"
	downloadURL := fmt.Sprintf("%s/api/books/download/%s", baseURL, token)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"download_url": downloadURL,
		"expires_in":   "1 hour", // Thời gian hết hạn
	})

	// Tự động xóa token sau 1 giờ
	time.AfterFunc(1*time.Hour, func() {
		bc.Lock()
		delete(bc.tempTokens, token)
		bc.Unlock()
	})
}

func (bc *BookController) DownloadFile(c *gin.Context) {
	token := c.Param("token")

	// Lấy bookID từ tempTokens
	bc.Lock()
	bookIDStr, exists := bc.tempTokens[token]
	if exists {
		delete(bc.tempTokens, token) // Xóa token sau khi dùng
	}
	bc.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Link tải không hợp lệ hoặc đã hết hạn",
		})
		return
	}

	// Chuyển đổi bookID sang uint
	var bookID uint
	if _, err := fmt.Sscanf(bookIDStr, "%d", &bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Lỗi hệ thống",
		})
		return
	}

	// Lấy thông tin sách từ database
	var book models.Book
	if err := bc.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Không tìm thấy sách",
		})
		return
	}
	log.Println("PDF URL from DB:", book.PDFUrl)

	relativePath := strings.TrimPrefix(book.PDFUrl, "/storage/")
	absolutePath := filepath.Join("/app/storage", relativePath)

	log.Println("Absolute path to file:", absolutePath)

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "File PDF không tồn tại",
			"debug":   absolutePath,
		})
		return
	}

	// Mở file
	file, err := os.Open(absolutePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Không thể mở file",
		})
		return
	}
	defer file.Close()

	// Lấy thông tin file
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Không thể đọc thông tin file",
		})
		return
	}

	// Thiết lập headers
	fileName := fmt.Sprintf("%s_%s.pdf", book.Title, book.Author)
	fileName = strings.ReplaceAll(fileName, " ", "_") // Thay thế khoảng trắng
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Gửi file
	if _, err := io.Copy(c.Writer, file); err != nil {
		log.Printf("Lỗi khi gửi file: %v", err)
	}
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

func SavePDFFile(file *multipart.FileHeader) (string, error) {
	uploadRoot := os.Getenv("UPLOAD_ROOT")
	if uploadRoot == "" {
		uploadRoot = "./storage" // fallback
	}
	log.Println("UPLOAD_ROOT:", uploadRoot)

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	savePath := filepath.Join(uploadRoot, "pdf", filename)
	log.Println("Save path:", savePath)

	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		log.Println("MkdirAll error:", err)
		return "", err
	}

	if err := saveUploadedFile(file, savePath); err != nil {
		log.Println("Save file error:", err)
		return "", err
	}

	log.Println("File saved successfully:", savePath)
	publicURL := "/storage/pdf/" + filename
	return publicURL, nil
}

func SaveImageFile(file *multipart.FileHeader) (string, error) {
	uploadRoot := os.Getenv("UPLOAD_ROOT")
	if uploadRoot == "" {
		uploadRoot = "./storage" // fallback
	}

	// Tạo tên file độc đáo với timestamp
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))

	// Đường dẫn lưu trong thư mục images
	savePath := filepath.Join(uploadRoot, "images", filename)

	// Tạo thư mục nếu chưa tồn tại
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Lưu file
	if err := saveUploadedFile(file, savePath); err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	// Trả về đường dẫn public
	publicURL := "/storage/images/" + filename
	return publicURL, nil
}

// Hàm save an toàn hơn cho gin
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
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

func (bc *BookController) GetKeywordsAndTOC(c *gin.Context) {
	bookID := c.Param("id")

	var book models.Book
	if err := bc.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sách"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"keywords":   book.Keywords,
			"toc_titles": book.TOCTitles,
		},
	})
}

func (bc *BookController) UpdateKeywordsAndTOC(c *gin.Context) {
	bookID := c.Param("id")

	var request struct {
		Keywords  []string `json:"keywords"`
		TOCTitles []string `json:"toc_titles"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	result := bc.DB.Model(&models.Book{}).Where("id = ?", bookID).Updates(map[string]interface{}{
		"keywords":   pq.StringArray(request.Keywords),
		"toc_titles": pq.StringArray(request.TOCTitles),
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cập nhật thất bại"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cập nhật thành công",
	})
}

func (bc *BookController) SearchByKeywords(c *gin.Context) {
	searchTerm := strings.TrimSpace(c.Query("q"))
	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu từ khóa tìm kiếm"})
		return
	}

	// Loại bỏ các từ khóa trống trong mảng keywords và toc_titles
	cleanSearchTerm := strings.ToLower(searchTerm)

	var books []models.Book
	err := bc.DB.Where(
		`(EXISTS (
            SELECT 1 FROM unnest(keywords) AS k 
            WHERE k <> '' AND LOWER(k) LIKE ?
        ) OR EXISTS (
            SELECT 1 FROM unnest(toc_titles) AS t 
            WHERE t <> '' AND LOWER(t) LIKE ?
        ))`,
		"%"+cleanSearchTerm+"%",
		"%"+cleanSearchTerm+"%",
	).Find(&books).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Lỗi tìm kiếm",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    books,
	})
}
