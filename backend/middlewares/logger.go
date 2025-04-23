package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware ghi log chi tiết các request
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ghi thời gian bắt đầu xử lý
		startTime := time.Now()

		// Xử lý request
		c.Next()

		// Tính thời gian xử lý
		latency := time.Since(startTime)

		// Ghi log
		log.Printf("[%s] %s | Status: %d | Latency: %v | Client: %s\n",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
		)

		// Ghi log lỗi nếu có
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("Error: %v\n", err)
			}
		}
	}
}
