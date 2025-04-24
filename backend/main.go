package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/models"
	"github.com/Poloni84Learning/ebook-store/routes"
	"github.com/Poloni84Learning/ebook-store/seeds"
	"github.com/Poloni84Learning/ebook-store/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load cấu hình từ file .env
	cfg := config.LoadConfig()

	// Khởi tạo kết nối database
	db := initDatabase(cfg)

	// Tự động migrate các model
	autoMigrate(db)

	// Seed dữ liệu mẫu (nếu cần)
	if cfg.DebugMode {
		seeds.SeedAll(db, cfg)
	}

	// Khởi tạo router
	router := routes.SetupRouter(db, cfg)

	// Khởi động GC dọn token hết hạn (chạy ngầm mỗi 10 phút)
	utils.StartTokenBlacklistGC(10 * time.Minute)

	// Chạy server
	runServer(router, cfg)
}

func initDatabase(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.SSLMode,
		cfg.TimeZone)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if !cfg.DebugMode {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Cấu hình connection pool
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(cfg.MaxDBConn)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established")
	return db
}

func autoMigrate(db *gorm.DB) {
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Book{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
		&models.SystemConfig{},
		&models.BookCombo{},
		&models.ComboItem{},
	}

	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to auto-migrate model: %v", err)
		}
	}
	log.Println("Auto migration completed")
}

func runServer(router *gin.Engine, cfg *config.Config) {
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)

	log.Printf("Starting server on %s\n", serverAddr)
	log.Printf("Database: %s@%s:%s/%s\n",
		cfg.DBUser,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
