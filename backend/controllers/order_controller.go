package controllers

import (
	"net/http"
	"time"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	DB     *gorm.DB
	Config *config.Config
}

type OrderInput struct {
	OrderItems    []OrderItemInput `json:"order_items" binding:"required,min=1"`
	PaymentMethod string           `json:"payment_method" binding:"required,oneof=Card COD BankTransfer"`
}

type OrderItemInput struct {
	BookID   uint `json:"book_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
}

type UserOrderUpdateInput struct {
	OrderItems    []OrderItemInput `json:"order_items" binding:"required,min=1"`
	PaymentMethod string           `json:"payment_method" binding:"required,oneof=Card COD BankTransfer"`
}

type StaffOrderUpdateInput struct {
	Status string `json:"status" binding:"required,oneof=pending processing completed canceled"`
}

func NewOrderController(db *gorm.DB, cfg *config.Config) *OrderController {
	go monitorOrders(db)
	return &OrderController{DB: db, Config: cfg}
}

// CreateOrder - Tạo đơn hàng mới
func (oc *OrderController) CreateOrder(c *gin.Context) {
	userID := c.GetUint("userID")
	var input OrderInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tính toán tổng giá và kiểm tra kho
	var total float64
	var orderItems []models.OrderItem

	for _, item := range input.OrderItems {
		var book models.Book
		if err := oc.DB.First(&book, item.BookID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
			return
		}

		if book.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock for book " + book.Title})
			return
		}

		total += book.Price * float64(item.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			BookID:   book.ID,
			Quantity: item.Quantity,
			Price:    book.Price,
		})
	}

	order := models.Order{
		UserID:        userID,
		TotalAmount:   total,
		Status:        "pending",
		OrderItems:    orderItems,
		PaymentMethod: input.PaymentMethod,
	}

	if err := oc.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	if err := oc.DB.
		Preload("User").
		Preload("OrderItems.Book").
		First(&order, order.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created order"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": order})
}

// UserUpdateOrder - Cập nhật đơn hàng (dành cho user)
func (oc *OrderController) UserUpdateOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetUint("userID")

	var input UserOrderUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := oc.DB.Preload("OrderItems").First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Kiểm tra quyền
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this order"})
		return
	}

	// Chỉ cho phép cập nhật khi đơn ở trạng thái pending
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only pending orders can be updated"})
		return
	}

	// Xóa items cũ và thêm items mới
	oc.DB.Where("order_id = ?", order.ID).Delete(&models.OrderItem{})

	var newItems []models.OrderItem
	var newTotal float64

	for _, item := range input.OrderItems {
		var book models.Book
		if err := oc.DB.First(&book, item.BookID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
			return
		}

		newItems = append(newItems, models.OrderItem{
			OrderID:  order.ID,
			BookID:   book.ID,
			Quantity: item.Quantity,
			Price:    book.Price,
		})
		newTotal += book.Price * float64(item.Quantity)
	}

	if err := oc.DB.Create(&newItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order items"})
		return
	}

	// Cập nhật thông tin đơn hàng
	order.OrderItems = newItems
	order.TotalAmount = newTotal
	order.PaymentMethod = input.PaymentMethod

	if err := oc.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	if err := oc.DB.
		Preload("User").
		Preload("OrderItems.Book").
		First(&order, order.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// StaffUpdateOrder - Cập nhật trạng thái đơn hàng (dành cho staff/admin)
func (oc *OrderController) StaffUpdateOrder(c *gin.Context) {
	orderID := c.Param("id")

	var input StaffOrderUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := oc.DB.
		Preload("User").
		Preload("OrderItems.Book").
		First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Cập nhật trạng thái
	order.Status = input.Status

	var updatedOrder models.Order
	if err := oc.DB.
		Preload("User").
		Preload("OrderItems.Book").
		First(&updatedOrder, orderID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetUserOrders lấy danh sách đơn hàng của user hiện tại
func (oc *OrderController) GetUserOrders(c *gin.Context) {
	userID := c.GetUint("userID")

	var orders []models.Order
	if err := oc.DB.
		Preload("User").            // Thêm User
		Preload("OrderItems.Book"). // Thêm Book trong OrderItems
		Where("user_id = ?", userID).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderDetails lấy chi tiết đơn hàng
func (oc *OrderController) GetOrderDetails(c *gin.Context) {
	userID := c.GetUint("userID")
	orderID := c.Param("id")

	var order models.Order
	query := oc.DB.
		Preload("User").           // Thêm User
		Preload("OrderItems.Book") // Thêm Book

	// Phân quyền
	if c.GetString("role") == string(models.RoleCustomer) {
		query = query.Where("id = ? AND user_id = ?", orderID, userID)
	} else {
		query = query.Where("id = ?", orderID)
	}

	if err := query.First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (oc *OrderController) GetAllOrders(c *gin.Context) {
	var orders []models.Order

	// Preload OrderItems và Book trong từng OrderItem
	if err := oc.DB.
		Preload("User").
		Preload("OrderItems.Book"). // Quan hệ lồng nhau
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách đơn hàng"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// monitorOrders - Tự động xử lý đơn hàng
func monitorOrders(db *gorm.DB) {
	for {
		time.Sleep(1 * time.Hour)

		// Hủy đơn chưa thanh toán sau 24h
		db.Model(&models.Order{}).
			Where("status = 'pending' AND created_at < ?", time.Now().Add(-24*time.Hour)).
			Update("status", "canceled")

		// Xóa đơn đã hủy sau 48h
		db.Where("status = 'canceled' AND updated_at < ?", time.Now().Add(-48*time.Hour)).
			Delete(&models.Order{})
	}
}
