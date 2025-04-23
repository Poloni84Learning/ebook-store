package controllers

import (
	"net/http"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComboController struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewComboController(db *gorm.DB, cfg *config.Config) *ComboController {
	return &ComboController{DB: db, Config: cfg}
}

// GetCombos - Lấy tất cả combos với phân trang
func (cc *ComboController) GetCombos(c *gin.Context) {
	var combos []models.BookCombo

	if err := cc.DB.
		Preload("User").
		Preload("ComboItems", "is_hidden = ?", false). // Chỉ lấy items không bị ẩn
		Preload("ComboItems.Book").
		Find(&combos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch combos"})
		return
	}

	c.JSON(http.StatusOK, combos)
}

// GetComboDetails - Lấy chi tiết combo
func (cc *ComboController) GetComboDetails(c *gin.Context) {
	id := c.Param("id")
	var combo models.BookCombo

	if err := cc.DB.
		Preload("User").
		Preload("ComboItems", "is_hidden = ?", false).
		Preload("ComboItems.Book").
		First(&combo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Combo not found"})
		return
	}

	c.JSON(http.StatusOK, combo)
}

// CreateCombo - Tạo combo mới
func (cc *ComboController) CreateCombo(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		BookIDs     []uint `json:"book_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra sách tồn tại
	var count int64
	if err := cc.DB.Model(&models.Book{}).Where("id IN ?", input.BookIDs).Count(&count).Error; err != nil || count != int64(len(input.BookIDs)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book IDs"})
		return
	}

	// Tạo transaction
	var combo models.BookCombo
	err := cc.DB.Transaction(func(tx *gorm.DB) error {
		combo = models.BookCombo{
			Title:       input.Title,
			Description: input.Description,
			CreatedBy:   userID,
		}

		if err := tx.Create(&combo).Error; err != nil {
			return err
		}

		// Thêm combo items
		var items []models.ComboItem
		for _, bookID := range input.BookIDs {
			items = append(items, models.ComboItem{
				ComboID:  combo.ID,
				BookID:   bookID,
				IsHidden: false,
			})
		}

		return tx.Create(&items).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create combo"})
		return
	}

	// Load lại toàn bộ thông tin kèm quan hệ
	var createdCombo models.BookCombo
	if err := cc.DB.
		Preload("User").
		Preload("ComboItems.Book").
		First(&createdCombo, combo.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created combo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Combo created successfully",
		"data":    createdCombo,
	})
}

// UpdateCombo - Cập nhật combo
func (cc *ComboController) UpdateCombo(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		BookIDs     []uint `json:"book_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lấy combo hiện tại
	var combo models.BookCombo
	if err := cc.DB.First(&combo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Combo not found"})
		return
	}

	// Kiểm tra quyền
	if combo.CreatedBy != userID && c.GetString("role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Transaction
	err := cc.DB.Transaction(func(tx *gorm.DB) error {
		// Cập nhật thông tin cơ bản
		if input.Title != "" {
			combo.Title = input.Title
		}
		if input.Description != "" {
			combo.Description = input.Description
		}

		if err := tx.Save(&combo).Error; err != nil {
			return err
		}

		// Cập nhật danh sách sách nếu có
		if input.BookIDs != nil {
			// Xóa items cũ
			if err := tx.Where("combo_id = ?", combo.ID).Delete(&models.ComboItem{}).Error; err != nil {
				return err
			}

			// Thêm items mới
			var items []models.ComboItem
			for _, bookID := range input.BookIDs {
				items = append(items, models.ComboItem{
					ComboID:  combo.ID,
					BookID:   bookID,
					IsHidden: false,
				})
			}
			return tx.Create(&items).Error
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update combo"})
		return
	}

	// Load lại toàn bộ thông tin
	var updatedCombo models.BookCombo
	if err := cc.DB.
		Preload("User").
		Preload("ComboItems.Book").
		First(&updatedCombo, combo.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated combo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Combo updated successfully",
		"data":    updatedCombo,
	})
}

// DeleteCombo - Xóa combo
func (cc *ComboController) DeleteCombo(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")

	var combo models.BookCombo
	if err := cc.DB.First(&combo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Combo not found"})
		return
	}

	// Kiểm tra quyền
	if combo.CreatedBy != userID && c.GetString("role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Xóa combo (ComboItems sẽ tự động xóa nhờ ON DELETE CASCADE)
	if err := cc.DB.Delete(&combo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete combo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Combo deleted successfully"})
}
