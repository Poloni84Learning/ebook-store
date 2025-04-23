package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	JWTExpiration time.Duration
	ServerPort    string
	SSLMode       string
	TimeZone      string
	MaxDBConn     int
	DebugMode     bool
}

func LoadConfig() *Config {
	return &Config{
		DBHost:        getEnv("DB_HOST", "db"), // Sử dụng "db" cho Docker Compose
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "ebook_store"),
		JWTSecret:     getEnv("JWT_SECRET", "default_secret_should_be_changed"),
		JWTExpiration: parseDuration(getEnv("JWT_EXPIRATION", "24h")),
		ServerPort:    getEnv("SERVER_PORT", "8081"),
		SSLMode:       getEnv("SSL_MODE", "disable"),
		TimeZone:      getEnv("TIME_ZONE", "Asia/Ho_Chi_Minh"),
		MaxDBConn:     parseInt(getEnv("MAX_DB_CONN", "10")),
		DebugMode:     parseBool(getEnv("DEBUG_MODE", "false")),
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseDuration(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 24 * time.Hour // Default to 24 hours if parsing fails
	}
	return duration
}

func parseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 10 // Default value
	}
	return val
}

func parseBool(s string) bool {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return false // Default value
	}
	return val
}
