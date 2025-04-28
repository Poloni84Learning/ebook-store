package controllers

import (
	"net/http"

	"github.com/Poloni84Learning/ebook-store/models"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type SystemConfigController struct {
	DB *gorm.DB
}

// CreateSystemConfig tạo mới SystemConfig
func (s *SystemConfigController) CreateSystemConfig(c *gin.Context) {
	var config models.SystemConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.DB.Create(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// GetSystemConfig lấy thông tin SystemConfig (giả sử chỉ có 1 bản ghi)
func (s *SystemConfigController) GetSystemConfig(c *gin.Context) {
	var config models.SystemConfig
	if err := s.DB.First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System config not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateSystemConfig cập nhật SystemConfig
func (s *SystemConfigController) UpdateSystemConfig(c *gin.Context) {
	var config models.SystemConfig
	if err := s.DB.First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System config not found"})
		return
	}

	var input models.SystemConfig
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cập nhật các trường
	config.ShippingFee = input.ShippingFee
	config.Promotion = input.Promotion
	config.PromotionInfo = input.PromotionInfo

	if err := s.DB.Save(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteSystemConfig xóa SystemConfig
func (s *SystemConfigController) DeleteSystemConfig(c *gin.Context) {
	var config models.SystemConfig
	if err := s.DB.First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "System config not found"})
		return
	}

	if err := s.DB.Delete(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "System config deleted successfully"})
}
