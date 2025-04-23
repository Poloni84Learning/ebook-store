package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/Poloni84Learning/ebook-store/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware tạo middleware xác thực JWT
func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bỏ qua các route public
		if isPublicRoute(c) {
			c.Next()
			return
		}

		// Lấy token từ header
		tokenString := utils.ExtractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}
		if utils.IsTokenBlacklisted(tokenString) {
			log.Printf("[DEBUG] Token bị từ chối vì nằm trong blacklist")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token đã bị vô hiệu hóa"})
			return
		}
		// Parse và validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			log.Printf("[DEBUG] Token parse error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if !token.Valid {
			log.Printf("[DEBUG] Token is invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Lấy claims và lưu vào context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			log.Printf("[TOKEN] Token claims: %#v", claims)
			setUserContext(c, claims, cfg)
		} else {
			log.Printf("[DEBUG] Token claims không đúng kiểu MapClaims, type = %T", token.Claims)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Next()
	}
}

// isPublicRoute kiểm tra route có public không
func isPublicRoute(c *gin.Context) bool {
	publicPaths := map[string]bool{
		"/api/healthcheck":   true,
		"/api/auth/register": true,
		"/api/auth/login":    true,
	}

	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/api/books/") && c.Request.Method == "GET" {
		return true
	}

	return publicPaths[path]
}

// setUserContext lưu thông tin user vào context
func setUserContext(c *gin.Context, claims jwt.MapClaims, cfg *config.Config) {
	// Parse user_id an toàn
	var userID uint
	if v, ok := claims["user_id"].(float64); ok {
		userID = uint(v) // Convert về uint chuẩn
	}

	// Lưu các thông tin cơ bản
	c.Set("userID", userID) // Đảm bảo userID luôn là uint
	c.Set("username", toString(claims["username"]))
	c.Set("role", strings.ToLower(toString(claims["role"])))

	// Thêm thông tin vào header cho debug (chỉ development)
	if cfg.DebugMode {
		c.Header("X-User-ID", toString(claims["user_id"]))
		c.Header("X-Username", toString(claims["username"]))
		c.Header("X-User-Role", toString(claims["role"]))
	}
}

// toString chuyển đổi interface{} sang string
func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case json.Number:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
