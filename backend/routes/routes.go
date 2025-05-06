package routes

import (
	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/controllers"
	"github.com/Poloni84Learning/ebook-store/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Middleware chung
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.LoggerMiddleware())
	router.Static("/uploads/covers", "/app/public/uploads/covers")

	router.Static("/storage", "/app/storage")

	// Khởi tạo controllers
	authController := controllers.NewAuthController(db, cfg)
	bookController := controllers.NewBookController(db, cfg)
	orderController := controllers.NewOrderController(db, cfg)
	comboController := controllers.NewComboController(db, cfg)
	reviewController := controllers.NewReviewController(db, cfg)
	systemConfigController := controllers.SystemConfigController{DB: db}

	// Public routes (không yêu cầu auth)
	public := router.Group("/api")
	{
		public.GET("/healthcheck", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "healthy"})
		})
		public.POST("/auth/register", authController.Register)
		public.POST("/auth/login", authController.Login)
		public.POST("/auth/staff-login", authController.StaffLogin) // New endpoint for staff/admin login
		public.GET("/books", bookController.GetBooks)
		public.GET("/books/:id", bookController.GetBook)
		public.GET("/books/by-title", bookController.GetBooksByTitle)
		public.GET("/books/by-author", bookController.GetBooksByAuthor)
		public.GET("/books/by-category", bookController.GetBooksByCategory)
		public.GET("/books/search", bookController.SearchBooks)
		public.GET("/combos", comboController.GetCombos)
		public.GET("/combos/:id", comboController.GetComboDetails)
		public.GET("/books/:id/combos", bookController.GetBookCombos)
		public.GET("/books/:id/reviews", reviewController.GetBookReviews) // Xem review sách
		public.GET("/books/top-selling", bookController.GetTopBooksByCompletedOrders)
		public.GET("/books/most-reviewed", reviewController.GetMostReviewedBooks)
		public.GET("/books/top-rated", reviewController.GetTopRatedBooks)
		public.GET("/categories", bookController.GetAllCategories)
	}

	// Protected routes (yêu cầu JWT auth)

	protected := router.Group("/api")
	protected.Use(middlewares.JWTAuthMiddleware(cfg))
	{
		protected.POST("/auth/logout", authController.Logout)
		// User routes
		user := protected.Group("/user")
		{
			user.GET("/profile", authController.GetProfile)
			user.PUT("/profile", authController.UpdateProfile)
		}

		// Book routes
		book := protected.Group("/books")
		{
			// Admin/Staff only
			adminBook := book.Group("").Use(middlewares.RoleMiddleware([]string{"admin", "staff"}))
			{
				adminBook.POST("", bookController.CreateBook)
				adminBook.PUT("/:id", bookController.UpdateBook)
				adminBook.DELETE("/:id", bookController.DeleteBook)

			}

			// Review routes
			book.POST("/:id/reviews", reviewController.CreateReview) // User đánh giá sách
			book.GET("/:id/download-link", bookController.GenerateDownloadLink)
			book.GET("/download/:token", bookController.DownloadFile)
			book.GET("/search-helper", bookController.SearchByKeywords)
			// // Review management
			review := protected.Group("/reviews")
			{
				review.PUT("/:id", reviewController.UpdateReview) // User cập nhật review của mình
				review.DELETE("/:id", reviewController.DeleteReview)
			}

		}

		// Order routes
		order := protected.Group("/orders")
		{
			order.POST("", orderController.CreateOrder)
			order.GET("", orderController.GetUserOrders)
			order.GET("/:id", orderController.GetOrderDetails)
			order.PUT("/:id", orderController.UserUpdateOrder) // User cập nhật đơn hàng

			// Admin/Staff only
			adminOrder := order.Group("").Use(middlewares.RoleMiddleware([]string{"admin", "staff"}))
			{
				adminOrder.GET("/all", orderController.GetAllOrders)
				adminOrder.PUT("/:id/status", orderController.StaffUpdateOrder) // Cập nhật trạng thái
			}
		}
		combo := protected.Group("/combos")
		{
			adminCombo := combo.Group("").Use(middlewares.RoleMiddleware([]string{"admin", "staff"}))
			{
				adminCombo.POST("/", comboController.CreateCombo)
				adminCombo.PUT("/:id", comboController.UpdateCombo)
				adminCombo.DELETE("/:id", comboController.DeleteCombo)
			}
		}

		// Admin only routes
		admin := protected.Group("/admin")
		admin.Use(middlewares.RoleMiddleware([]string{"admin"}))
		{
			admin.GET("/users", authController.ListUsers)
			admin.POST("/users", authController.CreateStaff)
			admin.PUT("/users/:id/role", authController.ChangeUserRole)
			admin.GET("/books/:id/keywords", bookController.GetKeywordsAndTOC)
			admin.PUT("/books/:id/keywords", bookController.UpdateKeywordsAndTOC)
			adminDashboard := admin.Group("/dashboard")
			{
				adminDashboard.GET("/top-books", bookController.GetTopBooks)
				adminDashboard.GET("/top-categories", bookController.GetTopCategories)
				adminDashboard.GET("/top-authors", bookController.GetTopAuthors)
				adminDashboard.GET("/top-total-orders", bookController.GetTotalOrders)
				adminDashboard.GET("/total-order", orderController.GetAllOrdersDashboard)
				adminDashboard.GET("/total-stats", orderController.GetOrderStats)
				adminDashboard.GET("/top-trending", reviewController.GetBookCountAboveRating)
				adminDashboard.GET("/order-trend", orderController.GetOrderTrends)
			}
			adminSystemConfig := admin.Group("/system-config")
			{
				adminSystemConfig.POST("", systemConfigController.CreateSystemConfig)
				adminSystemConfig.GET("", systemConfigController.GetSystemConfig)
				adminSystemConfig.PUT("", systemConfigController.UpdateSystemConfig)
				adminSystemConfig.DELETE("", systemConfigController.DeleteSystemConfig)
			}

		}
	}

	return router
}
