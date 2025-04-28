package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware tạo middleware kiểm tra role
func RoleMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		// Lấy role từ context (đã được set bởi JWTAuthMiddleware)
		userRole, exists := c.Get("role")
		if !exists {
			log.Println("[DEBUG] Role chưa set trong context")
		} else {
			log.Printf("[DEBUG] Role trong context (RoleMiddleware): %s", userRole)
		}
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Role information not found",
			})
			return
		}

		// Kiểm tra role có trong danh sách allowedRoles không
		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Invalid role type",
			})
			return
		}

		if !isRoleAllowed(roleStr, allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You don't have permission to access this resource",
			})
			return
		}

		c.Next()
	}
}

// isRoleAllowed kiểm tra role có được phép không
func isRoleAllowed(userRole string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if role == userRole {
			return true
		}
	}
	return false
}

// AdminOnlyMiddleware middleware chỉ cho phép admin
func AdminOnlyMiddleware() gin.HandlerFunc {
	return RoleMiddleware([]string{"admin"})
}

// StaffOnlyMiddleware middleware chỉ cho phép staff
func StaffOnlyMiddleware() gin.HandlerFunc {
	return RoleMiddleware([]string{"staff"})
}

// AdminOrStaffMiddleware middleware cho phép cả admin và staff
func AdminOrStaffMiddleware() gin.HandlerFunc {
	return RoleMiddleware([]string{"admin", "staff"})
}
