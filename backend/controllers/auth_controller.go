package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/models"
	"github.com/Poloni84Learning/ebook-store/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB     *gorm.DB
	Config *config.Config
}

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileInput struct {
	FirstName *string `json:"first_name" binding:"omitempty"`
	LastName  *string `json:"last_name" binding:"omitempty"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Phone     *string `json:"phone" binding:"omitempty"`
	Address   *string `json:"address" binding:"omitempty"`
	AvatarURL *string `json:"avatar_url" binding:"omitempty"`
}

type UserResponse struct {
	Success bool        `json:"success"`
	User    interface{} `json:"user"`
	Message string      `json:"message"`
}

func NewAuthController(db *gorm.DB, cfg *config.Config) *AuthController {
	return &AuthController{DB: db, Config: cfg}
}

func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput
	log.Println("[Auth] Bắt đầu xử lý đăng ký người dùng mới")
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[Auth] Lỗi validate input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Kiểm tra user đã tồn tại
	log.Printf("[Auth] Kiểm tra user đã tồn tại với username: %s hoặc email: %s\n", input.Username, input.Email)
	var existingUser models.User
	if err := ac.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		log.Printf("[Auth] User đã tồn tại: %+v\n", existingUser)
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "Username or email already in use",
		})
		return
	}

	// Hash password
	log.Println("[Auth] Bắt đầu hash password")
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		log.Printf("[Auth] Lỗi khi hash password: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Unable to create account",
		})
		return
	}

	// Tạo user mới

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         models.RoleCustomer,
		LastLogin:    nil,
	}
	log.Printf("[Auth] Bắt đầu tạo user mới: %+v\n", user)
	if err := ac.DB.Create(&user).Error; err != nil {
		log.Printf("[Auth] Lỗi khi tạo user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error creating account",
		})
		return
	}
	log.Printf("[Auth] Đăng ký thành công cho user ID: %d\n", user.ID)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registration successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (ac *AuthController) StaffLogin(c *gin.Context) {
	var input LoginInput
	log.Println("[Auth] Bắt đầu xử lý đăng nhập staff/admin")

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[Auth] Lỗi validate input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Tìm user trong database
	log.Printf("[Auth] Tìm kiếm staff/admin với username: %s\n", input.Username)
	var user models.User
	if err := ac.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		log.Printf("[Auth] Không tìm thấy user: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Incorrect login information",
		})
		return
	}

	// Kiểm tra role - chỉ cho phép staff hoặc admin
	if user.Role != models.RoleStaff && user.Role != models.RoleAdmin {
		log.Printf("[Auth] User không có quyền truy cập: Role=%s\n", user.Role)
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Account does not have access to the administrative system",
			"hint":    "Only staff and admin can log in here",
		})
		return
	}

	// Kiểm tra password
	log.Println("[Auth] Kiểm tra password hash cho staff/admin")
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		log.Println("[Auth] Password không khớp cho staff/admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect login information"})
		return
	}

	// Tạo JWT token với thời gian sống dài hơn (tuỳ chọn)
	log.Println("[Auth] Tạo JWT token cho staff/admin")
	token, err := utils.GenerateJWT(user.ID, user.Username, string(user.Role))
	if err != nil {
		log.Printf("[Auth] Lỗi khi tạo token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Unable to generate token",
		})
		return
	}

	// Cập nhật last login
	now := time.Now()
	log.Printf("[Auth] Cập nhật last login cho staff/admin: %v\n", now)
	ac.DB.Model(&user).Update("last_login", &now)

	log.Printf("[Auth] Đăng nhập thành công cho staff/admin ID: %d, Role: %s\n", user.ID, user.Role)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login to the management system successfully",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"is_admin": user.Role == models.RoleAdmin, // Thêm field đặc biệt cho admin
		},
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput
	log.Println("[Auth] Bắt đầu xử lý đăng nhập")
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[Auth] Lỗi validate input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Tìm user trong database
	log.Printf("[Auth] Tìm kiếm user với username: %s\n", input.Username)
	var user models.User
	if err := ac.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		log.Printf("[Auth] Không tìm thấy user: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Incorrect login information",
		})
		return
	}

	// Kiểm tra password
	log.Println("[Auth] Kiểm tra password hash")
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		log.Println("[Auth] Password không khớp")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Incorrect login information",
		})
		return
	}

	// Tạo JWT token
	log.Println("[Auth] Tạo JWT token")
	token, err := utils.GenerateJWT(user.ID, user.Username, string(user.Role))
	if err != nil {
		log.Printf("[Auth] Lỗi khi tạo token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Unable to generate token",
		})
		return
	}

	// Cập nhật last login
	now := time.Now()
	log.Printf("[Auth] Cập nhật last login: %v\n", now)
	ac.DB.Model(&user).Update("last_login", &now)
	log.Printf("[Auth] Đăng nhập thành công cho user ID: %d\n", user.ID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (ac *AuthController) Logout(c *gin.Context) {
	log.Println("[Auth] Bắt đầu xử lý logout")

	// Lấy token từ header
	tokenString := utils.ExtractToken(c)
	if tokenString == "" {
		log.Println("[Auth] Không tìm thấy token hợp lệ trong header hoặc cookie")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Token not provided",
		})
		return
	}

	// Thêm token vào danh sách đen (blacklist)
	err := utils.AddToBlacklist(tokenString, ac.Config.JWTExpiration)
	if err != nil {
		log.Printf("[Auth] Lỗi khi thêm token vào blacklist: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error while logging out",
		})
		return
	}
	log.Println("[Auth] Đăng xuất thành công")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

func (ac *AuthController) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	log.Printf("[Auth] Bắt đầu lấy profile cho user ID: %d\n", userID)
	var user models.User
	if err := ac.DB.First(&user, userID).Error; err != nil {

		log.Printf("[Auth] Không tìm thấy user ID %d: %v\n", userID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}
	log.Printf("[Auth] Trả về profile cho user ID: %d\n", userID)
	c.JSON(http.StatusOK, UserResponse{
		Success: true,
		User:    user.ToResponse(),
		Message: "Get profile successfully",
	})
}

func (ac *AuthController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	log.Printf("[Auth] Bắt đầu cập nhật profile cho user ID: %d\n", userID)

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[Auth] Lỗi validate input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := ac.DB.First(&user, userID).Error; err != nil {
		log.Printf("[Auth] Không tìm thấy user ID %d: %v\n", userID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	// Cập nhật thông tin
	log.Printf("[Auth] Cập nhật thông tin user từ %+v sang %+v\n", user, input)
	// Cập nhật các trường nếu chúng không nil
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Phone != nil {
		user.Phone = *input.Phone
	}
	if input.Address != nil {
		user.Address = *input.Address
	}
	if input.AvatarURL != nil {
		user.AvatarURL = *input.AvatarURL
	}

	if err := ac.DB.Save(&user).Error; err != nil {
		log.Printf("[Auth] Lỗi khi cập nhật user ID %d: %v\n", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update profile",
		})
		return
	}
	log.Printf("[Auth] Cập nhật thành công profile cho user ID: %d\n", userID)
	c.JSON(http.StatusOK, UserResponse{
		Success: true,
		User:    user.ToResponse(),
		Message: "Profile updated successfully",
	})
}

func (ac *AuthController) ListUsers(c *gin.Context) {
	var users []models.User
	if err := ac.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch users",
		})
		return
	}

	var response []models.UserResponse
	for _, user := range users {
		response = append(response, *user.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   response,
		"message": "Users fetched successfully",
	})
}

func (ac *AuthController) CreateStaff(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to create user"})
		return
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         models.RoleStaff,
	}

	if err := ac.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to create staff"})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		Success: true,
		User:    user.ToResponse(),
		Message: "Create staff successfully",
	})
}

func (ac *AuthController) ChangeUserRole(c *gin.Context) {
	userID := c.Param("id")
	var input struct {
		Role string `json:"role" binding:"required,oneof=admin staff customer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	var user models.User
	if err := ac.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "User not found"})
		return
	}

	user.Role = models.Role(input.Role)
	if err := ac.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Success: true,
		User:    user.ToResponse(),
		Message: "Change user role successfully",
	})
}
